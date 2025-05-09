#first we need to login, assume we already have login.json inside the request directory
# local test

test_environment=$1
if [ -z "$test_environment" ]; then
   echo "Environment can not be empty, should be either local or dev"
   exit 1
fi

if [ -z "$2" ]; then
   echo "File can not be empty, should be passed as a file"
   exit 1
fi

if [ -z "$3" ]; then
   echo "URL is empty, assuming default one,  http://localhost:9808/project/906d0fd569404c59956503985b330132/content"
   url="http://localhost:9808/project/906d0fd569404c59956503985b330132/content"
else
   url=$3
fi
base_url="http://localhost:9808"

if [[ $test_environment == "local" ]]; then
  login_url="http://localhost:9808/login"
fi

if [[ $test_environment == "dev" ]]; then
  login_url="https://app.cerex.io/mes/login"
fi


curl  -X POST  -F "file=@$2" $url




