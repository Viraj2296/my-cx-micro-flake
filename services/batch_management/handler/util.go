package handler

import (
	"bytes"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/batch_management/handler/database"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	qrcode "github.com/skip2/go-qrcode"
	"image/png"
	"io/ioutil"
	"strconv"
	"strings"
)

func generateQRCodeLabel(bardCodeText string, printerSettingField *database.PrinterInfo) (error, string) {

	data, err := ioutil.ReadFile("../batch_management/resources/qr_code_template.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return err, ""
	}

	// Convert the byte slice to a string
	content := string(data)
	content = strings.Replace(content, "[COR_X]", strconv.Itoa(printerSettingField.CoordinateX), -1)
	content = strings.Replace(content, "[COR_Y]", strconv.Itoa(printerSettingField.CoordinateY), -1)
	content = strings.Replace(content, "[QR_CODE_TEXT]", bardCodeText, -1)
	//TODO remove it later after tested this, before move to prod, or replace this with logger.
	fmt.Println("generated label content", content)
	return nil, util.ToEncodeBase64String(content)
}
func generateBarCodeLabel(bardCodeText string, printerSettingField *database.PrinterInfo) (error, string) {

	data, err := ioutil.ReadFile("../batch_management/resources/raw_material_label_format.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return err, ""
	}

	// Convert the byte slice to a string
	content := string(data)
	//content = strings.Replace(content, "[LABEL_WIDTH]", strconv.Itoa(printerSettingField.LabelWidth), -1)
	//content = strings.Replace(content, "[LABEL_LENGTH]", strconv.Itoa(printerSettingField.LabelLength), -1)
	//content = strings.Replace(content, "[WIDTH]", strconv.Itoa(printerSettingField.BarcodeWidth), -1)
	//content = strings.Replace(content, "[SPACE]", strconv.Itoa(printerSettingField.BarcodeSpace), -1)
	//content = strings.Replace(content, "[HEIGHT]", strconv.Itoa(printerSettingField.BarcodeHeight), -1)
	content = strings.Replace(content, "[COR_X]", strconv.Itoa(printerSettingField.CoordinateX), -1)
	content = strings.Replace(content, "[COR_Y]", strconv.Itoa(printerSettingField.CoordinateY), -1)
	content = strings.Replace(content, "[BARCODE_TEXT]", bardCodeText, -1)
	//TODO remove it later after tested this, before move to prod, or replace this with logger.
	fmt.Println("generated label content", content)
	return nil, util.ToEncodeBase64String(content)
}

// this will generate the png encoded image, so front-end should use this follwing.
// <img [src]="'data:image/png;base64,' + base64Image" alt="Base64 Image">
func generateBarcode(barCode string, width int, height int) (error, string) {
	enc := oned.NewCode128Writer()
	fmt.Println("generating label using width:", width, " height :", height, " and barcode: ", barCode)
	img, _ := enc.Encode(barCode, gozxing.BarcodeFormat_CODE_128, width, height, nil)
	if img == nil {
		return errors.New("invalid Barcode"), ""
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println("generated label:", base64Str)
	return nil, base64Str
}

func generateQRCode(qrCodeText string) (error, string) {
	img, err := qrcode.Encode(qrCodeText, qrcode.Medium, 256)
	if err != nil {
		return errors.New("invalid Barcode"), ""
	}

	// Decode the byte slice to an image
	qrImage, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return errors.New("failed to decode QR code to image"), ""
	}

	// Encode the image to PNG format
	var buf bytes.Buffer
	if err := png.Encode(&buf, qrImage); err != nil {
		return errors.New("failed to encode image to PNG"), ""
	}

	// Convert the PNG bytes to a base64 string
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println("generated label:", base64Str)
	return nil, base64Str
}
