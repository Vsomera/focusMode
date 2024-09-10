
# Focus Mode ??

Command Line Tool for blocking distracting websites, helping you focus on what actually matters.
- Only Windows 10 & 11 are supported

## Installation

Install focus mode with go

```bash
go install github.com/Vsomera/focusMode
```




    
## Usage and Examples

```
PS C:\User> focusmode
    ___       ___       ___       ___       ___         ___       ___       ___       ___
   /\  \     /\  \     /\  \     /\__\     /\  \       /\__\     /\  \     /\  \     /\  \
  /::\  \   /::\  \   /::\  \   /:/ _/_   /::\  \     /::L_L_   /::\  \   /::\  \   /::\  \
 /::\:\__\ /:/\:\__\ /:/\:\__\ /:/_/\__\ /\:\:\__\   /:/L:\__\ /:/\:\__\ /:/\:\__\ /::\:\__\
 \/\:\/__/ \:\/:/  / \:\ \/__/ \:\/:/  / \:\:\/__/   \/_/:/  / \:\/:/  / \:\/:/  / \:\:\/  /
    \/__/   \::/  /   \:\__\    \::/  /   \::/  /      /:/  /   \::/  /   \::/  /   \:\/  /
             \/__/     \/__/     \/__/     \/__/       \/__/     \/__/     \/__/     \/__/

| Cli Tool to block distracting websites, run "focusmode help" for command info.

```
- Adding domain(s) to blacklist `focusmode add "www.example.com" ...`
```
PS C:\User> focusmode add "www.instagram.com" "www.steam.com"

Added domain(s) to Blacklist:

|  1 www.instagram.com
|  2 www.steam.com
```
- Listing domains `focusmode list`
```
PS C:\Users> focusmode list

Blacklist:

|  1 www.instagram.com
|  2 www.steam.com
```

- Removing a single domain `focusmode clean --d "www.example.com"`

```
PS C:\User> focusmode clean --d "www.instagram.com"

Remove www.instagram.com from blacklist?

|  Type 'y' to confirm [y/n] y
|  removed www.instagram.com from blacklist
```
- Clearing blacklist `focusmode clean`
```
PS C:\User> focusmode clean

Clear all domains?

|  Type 'y' to confirm [y/n] y
|  cleared all domains

```
