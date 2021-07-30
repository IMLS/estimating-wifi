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

If you have never used one of these tiny adapters, this video shows you how to insert the microSD card. When working with these, always be gentle; don't force anything. The cards are directional -- they will only go into the adapter or Raspberry Pi *one way*. Therefore, if you can't get the card in, and it goes [kattywumpus](http://www.todayifoundout.com/index.php/2014/04/origins-kitty-corner-catawampus-cat-words/), then stop, back it out, and try it another way.

<iframe width="560" height="315" src="https://www.youtube.com/embed/YDIWEOL2GYU" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Install the operating system

Now, open the RPi Imager. If you are on Windows, you'll... click somewhere. Possibly the "start menu?" (Does Windows still have a start menu?) If you're on a Mac, it will be wherever you installed it... probably your Applications folder. 

Once you have launched the Raspberry Pi Imager, click **Operating System**. Then, click **Choose** in order to select an operating system. You will want to select **Raspberry Pi OS (other)**, which will present you with a second menu. From the second menu, select **Raspberry Pi OS Full (32-bit)**. In other words, you want the full, 32-bit Raspberry Pi OS.

Next, click **Storage**. Again, you'll get a list of options. Select your microSD card. It will show up as some kind of *32GB Generic Blah Blah* or something. (We're aware that "or something" conceals a great many sins. However, we're uncertain exactly how it will appear on your computer.)

After selecting the operating system and the microSD card, click **Install**.

This will take a few minutes. It will install the Raspberry Pi OS to the card, and then verify the install. 

If the concept of an "operating system" is relatively new to you (or you just never really thought about it before), you can watch this 3m video while you're waiting for the install to finish. Or, you could sip coffee and stare into space. Both are just fine.

<iframe width="560" height="315" src="https://www.youtube.com/embed/fkGCLIQx1MI" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Exit the imager

Now, you can exit the imager, remove the microSD card, and get ready for the rest of the process. You do not need to eject the card safely or anything of the sort. Just pop it out, and head back to the [setup instructions]({{ "/setup/#continuehere" | prepend: site.baseurl }}).
