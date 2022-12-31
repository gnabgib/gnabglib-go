package test

import (
	"strings"
	"testing"
)

//Compares two strings, creates errors if they differ in length, and shows
// the values side by side with style markers to show the issue
func StringCompare(t *testing.T, found,expect string) {
	if found == expect {
		return
	}
	if len(found)!=len(expect) {
		t.Errorf("Expected %d chars, got %d",len(expect),len(found))
	}

	eArr:=strings.Split(expect,"")
	eLine:=make([]byte,len(expect)+8)//4 extra chars for style/unstyle stages
	copy(eLine,expect)
	marker:=make([]byte,len(expect)+1)//1 extra char if found is longer
	var i int
	var r rune
	for i,r=range found {
		marker[i]='-'
		if i>=len(eArr) {
			break
		}
		//If we got a tab, add it to the marker (visual accuracy)
		switch {
		case r=='\x7f':
			//Duplicate del in marker
			marker[i]='\x7f'
		case int(r)<32:
			//If the char is control (other than above), copy it to output
			marker[i]=byte(r)
		}

		//We known eErr[i] is just a rune (because split by empty does this), so can directly compare
		if string(r)!=eArr[i] {
			marker[i]='^'
			//Move the data up 4 places
			copy(eLine[i+4:],eLine[i:])
			//Change the style
			copy(eLine[i:],[]byte("\x1b[7m"))
			//Reset the style at the end
			copy(eLine[len(eLine)-4:],[]byte("\x1b[0m"))
			
			//Flag the diff being mid-string (no need to append a ^ at the end)
			i=-1
			break
		}
		//eLine[i]=eArr[i]
	}
	if i>=0 {
		//If i isn't -1 we didn't find a diff yet, so mark the end as the issue
		marker[i]='^'
	}
	
	t.Errorf("Expect, found:\n\"%s\"\n\"%s\"\n %s",eLine,found,marker);
}