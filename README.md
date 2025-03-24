## What's For Lunch???

On execution, polls the school API and parses out the menu for the day.
Sends the result as an SMS.

The idea is to schedule the binary via cron to customize when the text is received.

For texting, we'll keep it simple and send it through email via carrier SMS gateway

### Each carrier has a specific SMS gateway domain.

- AT&T: @txt.att.net
- Verizon: @vtext.com
- T-Mobile: @tmomail.net
- Boost Mobile: @sms.myboostmobile.com
- Cricket Wireless: @sms.cricketwireless.net

### execution

The bin will take the email as a parameter as well as the gmail application key for sending email.

`lunch foo@txt.att.net abcdefg`

the intent is to use this with a cron that runs M-F like this:

`0 6 * * 1-5 lunch foo@txt.att.net abcdefg`

...assuming lunch is a bin in the path
