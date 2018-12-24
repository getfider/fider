 # Single Image

 if remove == true, then remove set current bkey to "" on that resource
 if upload != nil, then upload file to blob storage and set new bkey to that resource
 
# Multiple images

if item.length == 0, do nothing
if item.length > 0 {
  for each item {
    if item.remove == true, then remove the attachment with item.bkey from that resource (if exists on that resource)
    if item.upload != nil, then upload file to blob storage and add the new bkey to that resource
  }
}

