package inform

import (
	// Standard library
	"text/template"
)

var templateOptionList = parseTemplate("option-list", `
The options currently set for '{{.ID}}' are:
> Prefix: {{with .Options.Prefix}}'{{.}}'{{else}}(None Set){{end}}
Change these options with 'option set <key> <value>'.`)

var templateStoryList = parseTemplate("story-list", `
{{if .Stories}}
The list of active stories for '{{.ID}}' are:
{{- range .Stories}}
> Name: '{{.Name}}'
> Created at: {{.CreatedAt.Format "Mon, 02 Jan 2006 15:04"}}
> Last updated at: {{.CreatedAt.Format "Mon, 02 Jan 2006 15:04"}}
{{end}}
{{else}}
There are currently no active stories available for '{{.ID}}'.
Add a new one with 'story add' or get more information with 'help story' and 'help story add'.
{{end}}`)

var templateWelcome = parseTemplate("welcome", `
Hi! 👋

It seems like this is the first time we've exchanged messages (if not, feel free to skip this), and might need some help getting up-to-speed with what I do. With no further ado:

My name is InformBot, and I'm a chat interface for Inform 7, a system for creating works of interactive fiction, via a natural-language interface. What this means, essentially, is that you can write interactive books (and more) in very much the same language that you use reading them (that is, if you tend to read books written in the English language).

Inform is much, much more useful than just for writing interactive fiction stories for individuals — it can be used to implement interactive interfaces of any kind, and keep track of complex worlds with complex state.

That's where I come in, and can help in both defining the rules and performing actions against them, both in direct messages and in group-chats.

For more information on how you can define these rules, and how to converse with me in group-chats, type 'help' and I'll respond with a list of topics you can look further into.`)

var templateHelp = parseTemplate("help", `
Inform itself is a large, complicated system, and help on writing rules way beyond the scope of this help text. You are in luck though, as Inform comes with a large amount of documentation on its website: http://inform7.com/doc

Feel free to ask any questions about, or report any issues with InformBot itself here: https://github.com/deuill/informbot`)

var templateUnknownCommand = parseTemplate("unknown-command", `
I don't understand what '{{.}}' means. Type 'help' for an overview of common commands.`)

var messageInvalidSession = `
I couldn't start the story successfully — %s.`

var messageStartedSession = `
Story '%s' successfully started.
Any subsequent meta-commands will have to be given a prefix (currently set to '%s'), and you can end this session by using the 'story end' command. Have fun! 🎉`

var messageAddedStory = `
Story '%s' successfully added to active list.`

var messageRemovedStory = `
Story '%s' successfully removed from active list.`

// TODO FIX THESE TO BE MORE GENERIC

var messageUnknownStory = `
You need to pass in both the story name and URL, e.g. 'story add some-name https://example.com/story.ni'.
Story names need to be one word (though they can contain hyphens or underscores), and not contain any spaces or other white-space characters.`

var messageInvalidStory = `
I couldn't add the story successfully — %s.`

var messageSetOption = `
Option '%s' successfully set to '%s'.`

var messageUnknownOption = `
You need to pass in both the option name and value, e.g. 'option set Prefix ?'.`

var messageInvalidOption = `
I couldn't set that option successfully — %s.`

var messageRunError = `
I could't run that command — %s.`

var messageUnknownError = `
Oops, something went wrong and I was unable to complete that request, give me a moment and try again (or ask whoever set me up for some help).`

func parseTemplate(name, content string) *template.Template {
	return template.Must(template.New(name).Parse(content))
}
