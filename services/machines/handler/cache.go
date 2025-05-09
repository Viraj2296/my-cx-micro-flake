package handler

import (
	"gorm.io/gorm"
	"strconv"
)

type Node struct {
	value Message
	next  *Node
}

type MessageList struct {
	head *Node
	tail *Node
	size int
}

func (v *MessageList) AddToFront(value Message) {
	newNode := &Node{
		value: value,
		next:  v.head,
	}
	v.head = newNode
	v.size++
}

// Function to add an element to the back of the linked list
func (v *MessageList) AddToBack(value Message) {
	newNode := &Node{
		value: value,
		next:  nil,
	}
	if v.head == nil {
		// The linked list is currently empty
		v.head = newNode
		v.tail = newNode
	} else {
		// Append the new node to the tail
		v.tail.next = newNode
		v.tail = newNode
	}

	v.size++
}

// Function to search for messages with a timestamp greater than the given value
func (v *MessageList) GetMessagesByTimestampGreaterThan(ts int64) []Message {
	var messages []Message

	current := v.head
	for current != nil {
		if current.value.TS > ts {
			messages = append(messages, current.value)
		}
		current = current.next
	}

	return messages
}

// Function to search for messages with a timestamp less than the given value
func (v *MessageList) GetMessagesByTimestampLessThan(ts int64) []Message {
	var messages []Message

	current := v.head
	for current != nil {
		if current.value.TS < ts {
			messages = append(messages, current.value)
		}
		current = current.next
	}

	return messages
}

// Function to search for messages with timestamps between the given range (exclusive of boundaries)
func (v *MessageList) GetMessagesByTimestampBetween(startTS, endTS int64) []Message {
	var messages []Message

	current := v.head
	for current != nil {
		if current.value.TS > startTS && current.value.TS < endTS {
			messages = append(messages, current.value)
		}
		current = current.next
	}

	return messages
}

type Cache struct {
	LastOrderedId        int32
	DatabaseConnection   *gorm.DB
	MessageTopicMessages map[string]MessageList
}

/*
topic_1 : message1, message 2, message 3....
topic_2 : message1, message 2, message 3
*/
func (v *Cache) loadMessages() {
	messageQuery := "select * from message order by desc id"
	var messages []Message
	// loading messages largest id first
	v.DatabaseConnection.Raw(messageQuery).Scan(&messages)
	if len(messages) > 0 {
		v.LastOrderedId = messages[0].Id
		// adding in to list
		for _, message := range messages {
			topicBasedMessageList := v.MessageTopicMessages[message.Topic]
			topicBasedMessageList.AddToBack(message)
			v.MessageTopicMessages[message.Topic] = topicBasedMessageList
		}
	} else {
		v.LastOrderedId = 0
	}
}

func (v *Cache) loadDelta(dbConnection *gorm.DB) {
	deltaQuery := "select * from message where id > " + strconv.Itoa(int(v.LastOrderedId))
	var messages []Message
	// loading messages largest id first
	dbConnection.Raw(deltaQuery).Scan(&messages)
	if len(messages) > 0 {
		v.LastOrderedId = messages[0].Id
		for _, message := range messages {
			topicBasedMessageList := v.MessageTopicMessages[message.Topic]
			topicBasedMessageList.AddToBack(message)
			v.MessageTopicMessages[message.Topic] = topicBasedMessageList
		}
	}
}

func (v *Cache) getMessagesForMachinesBetween(topicName string, startTime DateTimeInfo, endTime DateTimeInfo) []Message {

	listOfMessages := v.MessageTopicMessages[topicName]
	selectedMessages := listOfMessages.GetMessagesByTimestampBetween(startTime.DateTimeEpochMilli, endTime.DateTimeEpochMilli)
	return selectedMessages
}
