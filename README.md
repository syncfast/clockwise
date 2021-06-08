# &orarr; Clockwise

Clockwise is a meeting cost calculator designed to encourage more efficient
meetings. 

> The meetings will continue until morale improves.

## Installation
```bash
go get github.com/syncfast/clockwise
```

## Usage
```
$ clockwise
Clockwise is a meeting cost calculator designed to encourage more efficient meetings.

Usage:
  clockwise [flags]
  clockwise [command]

Available Commands:
  help        Help about any command
  run         Run clockwise
  set         Set the average annual salary of meeting participants
  version     Print version

Flags:
  -d, --debug   verbose logging
  -h, --help    help for clockwise

Use "clockwise [command] --help" for more information about a command.
```

## Configuration
Clockwise relies on a single configuration item: average annual salary per
participant. This can be approximate, but you can use resources like indeed.com,
glassdoor.com, and levels.fyi to be methodical in your approximation.

`averageSalary` defaults to $150,000 and can be modified by editing the
configuration file located at `~/.config/clockwise/clockwise.yaml` or by using
the `set` subcommand:

```bash
$ clockwise set
? Set average annual salary of meeting participants: (150000) 
```

## Input
Clockwise supports the ability to automatically scrape participant count from a
specified zoom meeting by passing in the appropriate zoom URL. Just run:
```bash
clockwise run --url <zoom meeting url>
```

Alternatively, you can leverage manual input mode, which collects user input via
a TUI (terminal UI) to manage participant count:
```bash
clockwise run --manual
```

## Output
### TUI
Clockwise outputs total cost to a TUI. It's refreshed every 500ms.

### Virtual camera
Clockwise can integrate with [OBS](https://obsproject.com/) to output the total
cost of a meeting as it increases in real time via a virtual camera that can be
used as a video source to passive aggressively remind everyone in the meeting
how much the meeting is costing. Anecdotally, raising awareness this way led to
a rapid reduction in meeting frequency and duration. To change something, track
it. 

## OBS Configuration
OBS configuration is a bit involved, but it's something that you only need to do
once. Long term, I would like to replace OBS with something like FFMPEG that
doesn't depend on an external GUI.

- Install [OBS](https://obsproject.com/).
    ```bash
    brew install --cask obs
    ```
- Launch OBS.
- In the Sources window, click the `+` icon. 
- Select `Text (FreeType 2)`.
- Click `OK`.
- Check `Read from file`.
- Click `Select Font`.
- Change the `Size` to `144`.
- Click `OK`.
- Scroll down to the `Text File (UTF-8 or UTF-16)` field and click the `Browse`
  button.
- Navigate to `~/Documents/clockwise`.
- Select `clockwise.txt` and click `OK`. 
- Feel free to modify the text font, color, and size to suit your preferences.
  The sky is the limit.
- In the lower right hand corner, click `Start Virtual Camera`. 
- You can optionally add your webcam as a video source as well and make the
  cost a text overlay.
  - In the Sources window, click the `+` icon.
  - Select `Video Capture Device`.
  - Click `OK`.
  - In the `Device` dropdown, select your webcam. 
  - Click `OK`.
  - In the Sources window, make sure that `Video Capture Device` is selected and
    either drag it under your text source or click the down arrow icon.
  - If necessary, stretch your video layer so that it fills the canvas.
  - You can toggle your webcam off at any point by clicking the eye icon to the
    right of `Video Capture Device` in the Sources window.

## Zoom configuration
- If your zoom client is open, restart it (to load the new virtual camera).
- In Zoom, in the lower left hand corner, select the small up arrow next to the
  `Start Video` button.
- Select `OBS Virtual Camera`.
- Click the `Start Video` button. 

That's it! Make sure to add snarky remarks when your meetings drag on for added
effect.

## To do
- [ ] Explore authenticated web scraping approach.
- [ ] Explore Zoom's captcha requirement.
- [ ] Document plugin flow to support other video conferencing platforms (Google
  Meet, Microsoft Teams, etc.).
- [ ] Unit test all the things.
- [ ] Explore replacing OBS with a solution that doesn't require a GUI, like
  ffmpeg. 
- [ ] Explore time-series based approach.
- [ ] Explore crash recovery solutions (pick up where you left off). 
- [ ] Handle browser startup more gracefully. 
- [ ] Make the TUI prettier. 

## Contributing
This project is in a very early stage of development. It is far from perfect.
Contributions are welcome and appreciated!