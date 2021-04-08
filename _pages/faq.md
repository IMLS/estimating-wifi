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

## What is logged?

Our approach involves watching to see what wifi devices are nearby. When we see a new device, we keep track of it (temporarily) on the Pi, but give it an anonymous name. "Phone:338," for example. By temporarily, we mean that if someone comes to the library in the morning, we will see their device as a new, unique device if they come back in the afternoon. 

Every minute, we report the devices that are nearby. We report the anonymized name ("Phone:338"), the time, a unique identifier for the Raspberry Pi that is providing the report, and some other "telemetry" about the Pi itself. (By "telemetry" we mean "things that help us understand how the Pi is behaving.") 

From this, it is possible to determine what *kinds* of wifi devices have been present, but not identify them uniquely. And, of course, we can estimate how many devices used a library's wifi, and for how long.

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