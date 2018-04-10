# DauerreFlag
Simple reverse proxy to serve images.


## Example

```ruby
# https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_120x44dp.png
Base64.urlsafe_encode64('https://www.google.com',padding:false) # "aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbQ"
open('http://127.0.0.1:8001/aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbQ/images/branding/googlelogo/2x/googlelogo_color_120x44dp.png')
```