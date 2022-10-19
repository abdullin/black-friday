


# 2022-10-19

Deployment plan for DDD Meetup.

- [ ] 2 Cores at 5.22 USD per month on Hetzner.

![](images/Pasted%20image%2020221019085200.png)

- [ ] Link to bitgn.com
- [ ] Throw in caddy in front
- [ ] Install Gitea behind the scenes. Expose UI as bitgn.com/gitea, behind a password
- [ ] Python Flask web App to:
	- Ask for a public and a username
	- Set a cookie and save in a text file
	- Refresh the website in background until response is set
	- When response is set - display the git clone url
- [ ] Golang for building

**Behind the scenes**:
- when somebody registers - send a telegram message with credentials
	- Create an account in gitea with username 
	- Add the provided public key
	- Clone the private repository
	- Add the build hook
	- Send the clone url
- Build hook
	- Runs the script
	- prints result to the console
	- Sends a copy to my telegram

## Defer
- Gitea integration with telegram