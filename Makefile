include $(GOROOT)/src/Make.inc

TARG=hamcrest
GOFILES=\
	core.go\
	matcher.go\
	
include $(GOROOT)/src/Make.pkg
