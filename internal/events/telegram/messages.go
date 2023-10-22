package telegram

const helpMessage = `
Hello! I am url-saver, a bot that helps you save and tag your links. I use machine learning to automatically generate tags based on the content of the links.

Here's how you can use me:
1. Simply send me a link that you want to save.
2. I will extract the title and content of the webpage and analyze them to generate tag that may match the link's content.
3. In the future, you can use these tags to quickly search for and filter your saved links.

Here are the available commands:
- /get: Get the first saved link and remove it from the list.
- /show_tags: Show all your tags.
- /show_all: Show all saved links.
- /remove: Remove a link from the list. Format: "/remove *link*"

If you have any questions or need help, simply type the command /help.

Happy saving!`

const (
	helloMessage          = "Hello!\n\n" + helpMessage
	NoSavedPagesMessage   = "No saved pages."
	SavedMessage          = "URL saved."
	alreadyExistsMessage  = "This URL is already saved."
	unknownCommandMessage = "Unknown command."
	pageRemovedMessage    = "Link successfully removed."
	noLinkMessage         = "No link in message."
	noTagsMessage         = "Your links have no tags."
	noURLsForTagMessage   = "You have no URLs for this tag."
)

/*
get - Get first saved link and remove it from list.
show_tags - Show all your tags.
show_all - Show all saved links.
remove - Remove link from list. Format: "/remove *link*"
*/
