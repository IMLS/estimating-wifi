---
title: Prepping the SD Card
layout: page
sidenav: false
---
{% assign counter = 0 %}

# {{page.title}}

This is just one part of setting up the Raspberry Pi for use in the WIFISESS data pilot. It needs to be done from a computer where you have administrative access. This is usually true of computers you own, and sometimes true of computers you use at work.

This process takes roughly 15 minutes. Ideally, you complete this process before you begin setting up your Raspberry Pi.

There are online videos that explain this as well. [Here's a 3 minute video](https://www.youtube.com/watch?v=l9WSup73KuI) that takes you through the process. You can search for "raspberry pi imager setup" or similar to find other/similar instructions and videos.

<iframe width="560" height="315" src="https://www.youtube.com/embed/l9WSup73KuI" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Download the Raspberry Pi Imager

The [Raspberry Pi Imager](https://www.raspberrypi.org/software/) is a tool provided by the Raspberry Pi Foundation to help in setting up RPis.

Download and install the imager. You may need administrative privileges on your machine in order to do so.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Plug in the microSD card

Find the microSD card and the adapter that came in the kit. Put the card in the adapter. Plug the adapter into a USB port on your computer or attached hub.

You may use your own adapter if you have one you are particularly fond of. The ones in the kits are rather small and fiddly.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Install the operating system

Now, open the RPi Imager. Under **Operating System**, click **Choose** in order to select an operating system. You will want to select FIXME.

Next, under **Storage**, click **Choose**. Select your microSD card. It will show up as some kind of *32GB Generic Blah Blah* or something. (We're aware that "or something" conceals a great many sins. However, we're uncertain exactly how it will appear on your computer.)

After selecting the operating system and the microSD card, click **Install**.

This will take a few minutes.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Exit the imager

Now, you can exit the imager, remove the microSD card, and get ready for the rest of the process.
