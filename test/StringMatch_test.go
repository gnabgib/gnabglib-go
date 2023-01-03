package test

import (
	"strings"
	"testing"
)

var firstMismatchTests=[]struct {
	expect,found string
	first int
}{
	{"Hi","Hello",1},
	{"Hi","",0},
	{"Hi there","Hi you",3},
	{"Hi\nthere","Hi\nyou",3},
	{"Hii","Hi",2},
	{"Hi","Hii",2},
	{"Hi\x7fthere", "Hi\x7fyou're",3},
}

func TestFirstMismatch(t *testing.T) {
	for _, rec := range firstMismatchTests {
		
		eRunes:=strings.Split(rec.expect,"")
		fRunes:=strings.Split(rec.found,"")
		n:=len(eRunes)
		if len(fRunes)>n {
			n=len(fRunes)
		}
		marker:=make([]byte,n+1)
		
		f:=findFirstMismatch(eRunes,fRunes,marker)
		if rec.first!=f {
			t.Errorf("Expecting %d, found %d",rec.first,f)
		}
	}
}

var lastMismatchTests=[]struct {
	expect,found string
	last int
}{
	{"Hi","Hello",0},
	{"Hi","",0},
	{"Hi there","Hi you",0},
	{"Hi\nthere","Hi\nyou",0},
	{"Hii","Hi",-1},
	{"Hi","Hii",-1},
	{"Hi\x7fthere", "Hi\x7fyou're",-2},
}

func TestLastMismatch(t *testing.T) {
	for _, rec := range lastMismatchTests {
		
		eRunes:=strings.Split(rec.expect,"")
		fRunes:=strings.Split(rec.found,"")
		n:=len(eRunes)
		if len(fRunes)>n {
			n=len(fRunes)
		}
		marker:=make([]byte,n+1)
		
		l:=findLastMismatch(eRunes,fRunes,marker,0)
		if rec.last!=l {
			t.Errorf("Expecting %d, found %d",rec.last,l)
		}
	}
}

var renderTextTests=[]struct {
	text string
	blank,start,end int
	expect string
}{
	{"Hello",0,2,-1,"He\x1b[7mll\x1b[0mo"},
	{"Hello",1,2,-1,"He\x1b[7mll\xc2\xb7\x1b[0mo"},
	{"Hello",0,0,0,"\x1b[7mHello\x1b[0m"},
	{"Hello",1,0,0,"\x1b[7mHello\xc2\xb7\x1b[0m"},
	{"Hello",1,0,-4,"\x1b[7mH\xc2\xb7\x1b[0mello"},
}

func TestRenderText(t *testing.T) {
	for _, rec := range renderTextTests {
		f:=renderText(rec.text,rec.blank,rec.start,rec.end)
		if f!=rec.expect {
			t.Errorf("Expecting, found:\n%s\n%s",rec.expect,f)
		}
	}
}