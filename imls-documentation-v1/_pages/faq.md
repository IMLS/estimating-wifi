---
title: FAQ
layout: page
sidenav: false
---

# Questions you might have

We can imagine a question or two you might have. As more questions are asked, we'll update this list. 

## Who are you again?

We're federal employees. We're a team of engineers (James Tranovich and Matt Jadud) at <a href="https://18f.gsa.gov">18F</a>, a small digital services unit in the GSA. We work to make government more open and people-centered. 

The project overall has a much larger cast. We work closely with our colleagues in the Public Benefits portfolio at 18F, the work is funded by <a href="https://10x.gsa.gov">10x</a>, and we are doing it in close partnership <a href="https://imls.gov">IMLS</a>, whose vision guides our work. 

As we enter the pilot, we can say that this project is also in partnership with SDCs, State Librarians, and the public library community at large. *You*, in a word.

## Is privacy preserved?

No PII (personally identifiable information) is logged as part of this project. We believe it is impossible to use the data collected to identify an individual.

# What are we measuring/collecting?

What are we measuring/collecting in terms of data? Every minute, we are "watching" for (unencrypted, broadcast) wifi packets and noting "who" they're from. What we report/store looks like this:

```
"id","event_id","session_id","fcfs_seq_id","device_tag","servertime","localtime","manufacturer_index","patron_index"
177349,1044,"7cb76bc7f9186aa3","CA0001-001","in the depths of Z'ha'dum","2021-04-30T16:04:41Z","2021-04-30T16:04:37Z",4,8
177353,1044,"7cb76bc7f9186aa3","CA0001-001","in the depths of Z'ha'dum","2021-04-30T16:04:41Z","2021-04-30T16:04:37Z",4,26
177485,2219,"b6c5d2069cf728d3","ME0064-001","in the basement","2021-04-30T16:06:19Z","2021-04-30T16:06:18Z",21,59
177487,2219,"b6c5d2069cf728d3","ME0064-001","in the basement","2021-04-30T16:06:19Z","2021-04-30T16:06:18Z",6,37
```

The `id`, `event_id`, and `session_id` are "metadata," unrelated to people. The `id` is a global index; the `event_id` is an index within a session, and the `session_id` is reset every time the Raspberry Pi reboots. The `fcfs_seq_id` and `device_tag` are values you will enter when you set up the device. The timestamps are automatic. 

The `manufacturer_index` is a number that we assign to each device manufacturer as we see new devices. If the first device we see is made by Apple, the Apple is manufacturer zero. If the first device we see is made by Samsung, then Samsung is manufacturer zero. And, if your Pi resets, then these numbers reset, too. Put simply, this lets us track how many devices by a given manufacturer are showing up (within a session), but it anonymizes the value within and across the datasets. 

The `patron_index` is the same. If you walk into the library, and you are carrying the 37th device we've seen so far, then you become "patron 37." If you leave for more than 2 hours, when you come back, you'll become "patron 392". This resets whenever the device resets. Again, this cannot be traced back to an individual, as there is nothing personally identifiable about the index.

# What does this tell us?

With this logging, we can do some nifty filtering and (accurate) estimating as to how many people used your wifi for how long. Do we want to filter out devices we saw for 2 minutes or less? We can do that! Do we want to only count devices that were present for 15 minutes or more? We can do that! Are we dictating policy or determining how the WIFISESS element will change in the future by saying these things? **ABSOLUTELY NOT!**

But, we will have real data to work with so that questions about what we can actually observe/measure, what we can meaningfully report from the data, and what is involved in collecting that information can be meaningfully evaluated by the community. 

## Will the results of this work be open?

Yes. The code, data, and process by which we engage in this work will all be public domain/freely licensed. This website, even.

## Can I keep the Raspberry Pi?

Sadly, no. We'll be providing a mailing label so you can send them back. We'd like to let you keep them, but we don't make the rules. 

## Will I be compensated?

No. The GSA has some policy on that, and sadly, we can't. We don't think that's right, but we don't have scissors big enough to cut that particular bit of red tape.

## My library wasn't selected to be part of the pilot. Can we still take part?

We have a limited number of devices to distribute for the pilot. Further, we are constrained by process and regulation as to how many can take part in the pilot. This is why we are (currently) limited to 9 libraries in the pilot.

We are working on getting approval to allow more libraries to take part. If we succeed, we will update this site with more information as to how you can use your own Raspberry Pi to join in.

{% comment %}
However, if you have a Raspberry Pi, and you have a compatible wifi adapter (or are willing to purchase one), we'd love to have you take part.

A Raspberry Pi kit with everything you need costs between $70 and $100 dollars. A wifi adapter that will work with the code that has been written for this pilot will cost roughly $30 dollars. 
{% endcomment %}




## My questions weren't answered here. Where can I ask my question?
 
If you have questions, we have a [friendly little form](https://forms.gle/qTkUmGEErUi6Wcrn7) where you can submit them. Thank you! Matt or James will do their best to respond in a timely manner. We're a small team, so please be patient with us.