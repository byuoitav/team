echo "circlevars: updates environment variables in circleci. Please make sure your CIRCLE_API_TOKEN is set."

echo "Enter username:" && echo

read USERNAME

echo "Enter the name of the project to update:" && echo

read PROJECT

echo "Enter the name of the environment variable to update:" && echo 

read VARIABLE

curl -X DELETE https://circleci.com/api/v1.1/project/github/$USERNAME/$PROJECT/envvar/$VARIABLE?circle-token=$CIRCLE_API_TOKEN

curl https://circleci.com/api/v1.1/project/github/$USERNAME/$PROJECT/envvar?circle-token=$CIRCLE_API_TOKEN

echo "Enter new value for variable:" && echo

read NEW

curl -X POST --header "Content-Type: application/json" -d '{"name":"'$VARIABLE'", "value":"'$NEW'"}' https://circleci.com/api/v1.1/project/github/$USERNAME/$PROJECT/envvar?circle-token=$CIRCLE_API_TOKEN
