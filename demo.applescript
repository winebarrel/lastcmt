#!osascript
# NOTE: Changing the title of the target iTerm window

tell application "iTerm"
  set winlist to every window
  repeat with win in winlist
    set the_title to name of win
    if the_title contains "lastcmt demo" then
      activate
      set index of win to 1
      tell current window
        tell current session of current tab
          write text "clear"
          delay 5
          write text "export GITHUB_OWNER=winebarrel"
          delay 1
          write text "export GITHUB_REPO=hello-world"
          delay 1
          write text "echo 'こんにちは' | lastcmt -n 12 -"
          delay 5
          write text "echo '今日は' | lastcmt -n 12 -"
          delay 5
          write text "gh pr comment -R winebarrel/hello-world 12 -b 'gh comment' > /dev/null"
          delay 5
          write text "echo 'こんばんは' | lastcmt -n 12 -"
        end tell
      end tell
    end if
  end repeat
end tell
