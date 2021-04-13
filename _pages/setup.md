---
title: Setting up your Pi
layout: page
sidenav: false
---

These instructions are for setting up a Raspberry Pi for participation in the IMLS/18F/10x WIFISESS pilot. We recommend reading all of the instructions first, and then following them step-by-step.

If the setup fails in any way, please [reach out to the 18F team]({{site.questionformurl}}), and we'll connect with you to help complete the setup.

## THINGS YOU WILL NEED

*To take part in the pilot, you should have already been contacted, and you should have confirmed that you have the following available.*

To do this, you will need:

1. A USB keyboard
2. A USB mouse
3. A monitor with HDMI
4. An open ethernet port near an electrical outlet

It is possible you could use a monitor with some other connector (DVI) if you have the correct adapters. While we are amazing, we cannot help you with this. If you do not have an HDMI monitor, we recommend digging around for one in your box of weird dongles (you have one of those, right?), asking your IT team, a tech-savvy friend, or an area 15-year-old who is some combination of eager and bored to help you.

Many monitors have more than one input, and you can switch between them. So it should *not* be the case that you need to buy another monitor; you might be able to use the one right in front of  you. That's what we do.

The "ethernet port near an electrical outlet" means you have a place to plug the RPi into both power and network. (We're going to start calling it the **RPi**, because it's shorter.) Ideally, this is not buried in a closet, but is somewhere near the circ desk (or somewhere that you have high foot traffic).

**If you don't have these pieces, you cannot proceed.**

## Unbox your Pi

You should have been mailed a Raspberry Pi kit, a USB wifi adapter, and a short ethernet cable.

## Assemble the Pi

The video below walks you through assembling the Raspberry Pi.

<iframe width="560" height="315" src="https://www.youtube.com/embed/7rcNjgVgc-I" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

Googling "canakit raspberry pi kit assembly" should turn up additional videos and resources. We are not documenting this process ourselves because it is well-documented elsewhere on the internet.

## Download the Raspberry Pi Imager

The [Raspberry Pi Imager](https://www.raspberrypi.org/software/) is a tool provided by the Raspberry Pi Foundation to help in setting up RPis.

Download and install the imager. You may need administrative privileges on your machine in order to do so.

## Plug in the microSD card

Find the microSD card and the adapter that came in the kit. Put the card in the adapter. Plug the adapter into a USB port on your computer or attached hub.

You may use your own adapter if you have one you are particularly fond of. The ones in the kits are rather small and fiddly.

## Install the operating system

Now, open the RPi Imager. Under **Operating System**, click **Choose** in order to select an operating system. You will want to select FIXME.

Next, under **Storage**, click **Choose**. Select your microSD card. It will show up as some kind of *32GB Generic Blah Blah* or something. (We're aware that "or something" conceals a great many sins. However, we're uncertain exactly how it will appear on your computer.)

After selecting the operating system and the microSD card, click **Install**.

This will take a few minutes.

## Exit the imager

Now, you can exit the imager, remove the microSD card, and get ready for the slow part of the process.

## Insert the card

Put the microSD card in the RPi. The slot is on the bottom of the pi. The gold-plated contacts on the card will be facing UP, and the logo on the card will be facing DOWN.

Also, it will only go in *one way*, so if it doesn't slot in easily, *do not force it*. Flip it around and try again if you meet resistance.

## Plug in the pi

Now it is time to plug everything in.

1. Plug one end of the ethernet cable into the RPi, and the other end into an ethernet jack/switch/router or similar.
2. Plug one end of the micro-HDMI cable into the RPi (the little end) and the other end into your monitor.
3. Plug your keyboard and mouse in.
4. Plug in the USB wifi adapter.
5. Plug the power adapter into the switch cable, and the other end of the switch cable into the RPi.

## Turn it on

Now that everything is plugged in, turn it on.

## Wait

Wait.
