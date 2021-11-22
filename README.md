# twilio-secret-santa
Use Twilio's SMS API to organize Secret Santa!

#### Usage
1. Clone this repo and modify the [list of participants](https://github.com/abedef/twilio-secret-santa/blob/5abebb8e8d3166900183f98b463fd6ab793ac036/santa.go#L30) to match your Secret Santa participants
2. Set your Twilio environment variables `TWILIO_TOKEN`, `TWILIO_NUMBER`, and `TWILIO_SID` (see https://www.twilio.com/docs/sms)
3. Run santa.go (`go run santa.go` if you have already set your environment variables, or `TWILIO_TOKEN=??? TWILIO_NUMBER=??? TWILIO_SID=??? go run santa.go` to define the environment variables and run the program in one Go (lol))

If your Twilio credentials and phone numbers are valid, each participant in your [list of participants](https://github.com/abedef/twilio-secret-santa/blob/5abebb8e8d3166900183f98b463fd6ab793ac036/santa.go#L30) will recieve a [text](https://github.com/abedef/twilio-secret-santa/blob/5abebb8e8d3166900183f98b463fd6ab793ac036/santa.go#L71) greeting them and informing them of their Secret Santa recipient. For example, if Rudolph is assigned as Dasher's Secret Santa, Rudolph would recieve a text that says `Hey Rudolph! You are assigned to Dasher for Secret Santa!`.
