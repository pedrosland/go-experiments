# Image Blur

This app will blur the provided image using a box blur which averages the channels on 9 pixels to produce the new value.

Parallelism can be adjusted:
 * 0 uses entirely sequential code
 * 1 uses parallelism of 1
 * \>1 uses that n goroutines assuming that your image is at least n pixels wide

This was ported from a similar project I was doing in Scala.

## Attribution

The images are provided with thanks to the following people:

**Bird with fish**  
Nigel  
https://www.flickr.com/photos/winnu/28441769481/

**The Go Gopher**  
Renee French  
http://reneefrench.blogspot.com/
