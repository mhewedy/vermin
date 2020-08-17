# vermin website

### Technical details:

* we use the for of `<link-command></link-command>` to link other commands from command pages.
The result will page `<a href="#" onclick="$('#list-COMMAND-list').click()">COMMAND</a>`
    
  example:`<link-stop><link-stop>` will generate: `<a href="#" onclick="$('#list-stop-list').click()">stop</a>`


* If we need to override the text of the generate link, then we supply the text as part of the `link-xxx` element:
`<link-command>TEXT HERE</link-command>`

  example:`<link-stop>Stop VM<link-stop>` will generate: `<a href="#" onclick="$('#list-stop-list').click()">Stop VM</a>`
