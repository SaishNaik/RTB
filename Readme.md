**Real Time Bidding**

Real-time bidding ( RTB) is a subcategory of programmatic media buying. It refers to the practice of buying and selling ads in real time on a per-impression basis in an instant auction. This is usually facilitated by a supply-side platform (SSP) or an ad exchange. 

**Project Details:**
Only SSP component implemented, all of the other services are dummy data providers

**To start the services,**

1.Run from the root folder of project, below commands 

a.`docker compose up --build`

b.Run the below command for mongoseeding<br>
    ` docker-compose exec -T db mongorestore --archive --gzip < dump.gz` 
.The seeding needs to be run only once as volumes will take care of the other runs


Total Services:<br>
1.SSP<br>
2.Pub<br>
3.DSP<br>
4.Mongo<br>

Inorder to see the app running,type `http://localhost:3001/` in browser.

**Pub details:**<br>
Pub will have a div `<div class="admaru-ads" data-width=43 data-height=43 data-adtype="Banner" data-adslotid="6458c5c947743c1a71fa9f7b">`
.Currently using adslot id for each adslot which will contain pubid and site info saved in ssp db
We go through all data params of div tag and compose a ad request and send to ssp


**SSP details:**

1.SSP will validate the params sent and create a bid request. Bidrequest will be created using adslot data such as site etc. We will find country
using ip2location and os from useragent. Device information can be got from publisher using data-attributes(not implemented currently)

2.Then it conducts a auction. In the auction, we find list of dsps
attached to pub and then send request to each of those dsps. Once we get the response from all dsp's, we select the maximum bid.
<br>We then try to fetch markup in 2 ways:<br> 
a)if it is sent through bid response<br>
b)if it is sent through win url, we substitute the macros and call win url<br>
Once we get the markup, we add our impression url and return the markup.
After returning the markup, we save bid request details in db

3.When the ad is displayed, the impression url is fired and we check if the bidrequest exists in db and if it does,
we update the profit and revenue and impression in db.

**Note:** since bidresponse from the dsp contains impression url (from my experience),hence we are not redirecting when 
receive the impression url. Also some ads can be clickable and that can be saved as well on our end, but that was not
the scope here.

**Dsp Service:**
Just set 2 routes for dsp1 and dsp2, which return bid responses

**Mongo:**
just a normal container. Also I have persisted the data as a volume




