

To manage dependencies, run the database and server, please make sure the run script has execution permissions: 

chmod +x dev_run.sh


Then finally to run the database docker image and server:

./dev_run.sh 


or to manually run, without using the prepared script: 

To resolve dependencies please run: 

go mod tidy 

Then to run:

docker compose up -d 
go run .

