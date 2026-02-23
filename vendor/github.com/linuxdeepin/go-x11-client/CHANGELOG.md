[0.4.1] 2019-04-03
*   fix: build failed

[0.4.0] 2019-03-27
*   feat: add util xcursor

[0.3.0] 2019-02-26
*   chore: use WriteBool and ReadBool
*   feat: xkb add more functions
*   chore: exts use new reader
*   chore: xproto event use new reader
*   feat: implement new Reader
*   feat: add ext xkb
*   feat: delete other error types, leaving only the Error type
*   chore: continue to replace uint64 with SeqNum
*   feat: use special seqNum to indicate error
*   fix: requestCheck do not unlock c.ioMu before return
*   chore: replace uint64 with SeqNum
*   feat: reduce memory usage for encoding requests

# [0.2.0] - 2018-11-23
*   fix: readSetup

# [0.1.0] - 2018-10-25
*   fix(wm/ewmh): get icon failed
*   feat: add WriteSelectionNotifyEvent
*   chore: add makefile for `sw_64`

# [0.0.4] - 2018-07-19
*   fix: requestCheck method
*   feat: add ext shm
*   feat: add ext render
*   feat: support ge generic event
*   feat: add ext input
*   fix: readMapNotifyEvent no set Window field
*   fix: conn read blocked when event chan full
*   feat(ext/randr): add more events
*   fix(ext/randr): name typo
*   feat: add ext xfixes
*   fead: handle conn read write error
*   feat: Conn add method IDUsedCount
*   fix: wrong call NewError
*   feat: add resource id allocator
*   feat: conn add atom cache
*   fix: ewmh _NET_SUPPORTING_WM_CHECK is not root property
*   feat: add ext randr
*   feat: add ext dpms and screensaver
*   chore: call log.Print if env DEBUG_X11_CLIENT = 1
*   chore(deb): set env DH_GOLANG_EXCLUDES
*   feat: handwriting encode decode part
*   fix: atom WM_CHANGE_STATE not found

# [0.0.3] - 2018-03-07
*   remove usr/bin
*   expand event chan buffer size
*   fix: Adapt lintian

