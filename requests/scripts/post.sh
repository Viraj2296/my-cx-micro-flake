#first we need to login, assume we already have login.json inside the request directory
# local test

test_environment=$1
if [ -z "$test_environment" ]; then
   echo "Environment can not be empty, should be either local or dev"
   exit 1
fi

if [ -z "$2" ]; then
   echo "Post body can not be empty, should be passed as a file"
   exit 1
fi

if [ -z "$3" ]; then
   echo "URL can not be empty, should be a valid URL"
   exit 1
fi
base_url="http://localhost:9808"

if [[ $test_environment == "local" ]]; then
  login_url="http://localhost:9808/login"
fi

if [[ $test_environment == "dev" ]]; then
  login_url="https://app.cerex.io/mes/login"
fi

token=$(curl --silent -X POST --data @../request_sigatures/login.json -H "Content-Type: application/json" $login_url | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "Authenticated [" $token "]"

curl  -X POST --data @$2 -H "Content-Type: application/json" -H "Authorization: Bearer $token" $3




