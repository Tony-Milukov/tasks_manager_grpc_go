This is my tasks manager backEnd project build with Go (grpc)
    
Technologies used:
1. go as language
2. jwt for auth
3. postgres as database (sql / psql)
   1. also implemented auto migrations 
4. docker
5. bcrypt

Functionality

                             Auth:

    1. register user
    2. login

                             Tasks:
    1. CreateTask
    2. DeleteTask
    3. UpdateTask
    8. GetTasksByFilter
    9. UnAssignTask
                             Statuses:
    1. GetAllStatuses
    2. UpdateStatus
    3. CreateStatus
    4. DeleteStatus






I used docker and the best practise of Go development

How to start the project:

    1. Install Docker

    2. Run Docker

    3. Fork this repo with git clone _ repoUrl _ 

    4. Go to the repo dir

    5. Open terminal

    6. Exec this command: "docker compose up"

    // if you want to test it
    7. install postman https://www.postman.com/downloads/

    8. open postman

    9. select grpc as protocol

    10. Import protoFile
        select "import proto file";
        select api.proto from /proto/proto/api.proto

    11. Select the needed method and click on "Send"

Thank you, enjoy!
    
   