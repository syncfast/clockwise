# Extending Clockwise

My hope is that end users who rely on video conferencing platforms other than
Zoom can contribute scraping implementations for their respective platforms. In
theory, it should be as simple as:
- Duplicate [zoom.go](zoom.go). 
- Modify it to target a different platform. 
- Add the appropriate flags to [run.go](../../cmd/run.go).

User contributions are welcome and appreciated. I'm happy to field any questions
that surface in the process.