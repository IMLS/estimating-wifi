---
title: Setting up your Pi
layout: page
sidenav: false
blockerstep: 10
---

# {{page.title}}

These instructions are for setting up a Raspberry Pi for participation in the IMLS/18F/10x WIFISESS pilot. We recommend reading all of the instructions first, and then following them step-by-step.

If the setup fails in any way, please [reach out to the 18F team]({{site.questionformurl}}), and we'll connect with you to help complete the setup.

{% assign counter = 0 %}
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Read through the setup instructions

To complete the setup **you'll need roughly 1 hour**, assuming you have everything at your fingertips. Some of those steps require you to have particular tools, pieces-parts, or connectivity available. Please read the instructions and make sure you have everything you need to hand.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Set up the microSD card

You can do this step at any time; that said, it will prevent you from moving past step {{page.blockerstep}} in these instructions, so you might as well do it in advance. 

We have a separate page that describes how to [set up the microSD card]({{ "preppingthesdcard/" | prepend: site.baseurl }}). It takes roughly 15 minutes, and you will need to do it on a computer or laptop where you have administrative access.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Request an api.data.gov key

We have [a separate page for this]({{ "requestkey/" | prepend: site.baseurl }}). This is an easy process, but may take a day or two in the background. Go ahead and do this in advance of anything else.
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Things you must have

There are a few things you **must** have for this to work.

1. **A Raspberry Pi**.<br>
    You should have received an RPi kit if you are part of the pilot.
2. **A USB wifi adapter**.<br>
    Again, this should have been part of your kit. 
3. **You will need a USB keyboard**.<br>
    Any USB keyboard should work. You only need it to set up the RPi. You could borrow it from another computer.
4. **You will need a USB mouse**. <br>
    Like the keyboard, you could borrow this from elsewhere.
5. **You will need [a monitor or TV with HDMI in](https://www.youtube.com/watch?v=Il64oLobp38)**.<br>
    Many TVs also work for this purpose; you may need to use your "source" or "input" button to select the correct device if you have more than one thing plugged into your TV or monitor.
6. <span style="color: red;">**You will need access to an ethernet port**</span>.<br>
    In theory, you *might* be able to set up the RPi using a wifi network. We strongly recommend against it.
7. **A computer where you have admin access**.<br>
    For setting up the microSD card, you will need a computer or laptop where you have admin access. Any Windows, Mac, or Linux computer will work. 
8. **Optional**. A phone, tablet, laptop, or other device that can access the internet.<br>
    We think you might want these directions open on a separate device while you are working on the RPi.

The ethernet port might be difficult. We included a 3' ethernet cable in the kit, but it may be too short to be useful. Any ethernet cable will work.

**If you don't have these pieces, you cannot proceed.**

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Find a place to do the setup

You can do this setup anywhere you like, as long as the conditions above are met. 

Because many of us are working from home/living at work, you may find that you're doing this at home. 

In my house, this means I would probably be doing the setup *in the living room*. Why? Following from the requirements above:

1. I can take the RPi to the living room.
2. The USB wifi can also be carried to the living room.
3. I can borrow a keyboard from another computer elsewhere in the house.
4. I can borrow a mouse from another computer in the house.
5. My TV has HDMI in. I may have to unplug the DVD player, but that's OK.
6. My wifi router (which is next to the TV) has 4 ethernet ports. I may have to unplug the Wii (because all four of the ports are currently in use), but no one will be using the Wii while I am using the TV to set up the RPi.
7. I have a laptop I can use to read these directions while I sit on the floor/couch and swear at whomever wrote these instructions.

You might rather do this setup at the library. Note the [possible points of failure](#possible-points-of-failure) at the end of this document. Point #1 highlights that *things might go wrong* if you do this in the library. We know enough about networks to know that *something might not work*, but we don't know *exactly what*. If you have network issues during setup, you will probably need to talk to your local IT team. If you do not have a local IT team, then please feel free to reach out to the 18F team, and we will do our best to see you through to success.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Ready yourself

If you are a coffee drinker, this is where you refresh your cup. If you are a tea drinker, get a fresh pot going. Get everything set, and perhaps have these directions up on your phone, tablet, or laptop off to the side.
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Unbox your Pi

You should have been mailed a Raspberry Pi kit, a USB wifi adapter, and a short ethernet cable. Unbox the RPi!

It will have a bunch of bits in it. Lay them out. Aren't they fun? We especially like the cute little fan.
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Assemble the Pi

The video below walks you through assembling the Raspberry Pi. At roughly 7:30m, you will see the assembly is complete, and the startup process begins. (Wagner goes on to run some commands and things on the RPi; we ask you do *not* do this.)

<iframe width="560" height="315" src="https://www.youtube.com/embed/7rcNjgVgc-I" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

We also found a [wonderful video by Karunesh](https://www.youtube.com/watch?v=deJQFAPfRZg) that talks you through the whole process as well. There are clearly haters on the internet, because 26 people gave Karunesh a thumbs down. We think the video does a great job of introducing the RPi and the setup process. The first few minutes are an overview, and [at roughly 3:30s, Karunesh and his trusty helper assemble the RPi](https://youtu.be/deJQFAPfRZg?t=216). 

<iframe width="560" height="315" src="https://www.youtube.com/embed/deJQFAPfRZg" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

Googling "canakit raspberry pi kit assembly" should turn up additional videos and resources. We are not documenting this process ourselves because it is well-documented elsewhere on the internet. 

If this step proves to be a challenge, please reach out to the 18F team. We can expand this documentation and will answer any questions you have.

When you are done, you will have stuck on the heatsinks and fan, inserted it into the case, and it will be sitting there, disconnected, unpowered, but ready to go.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Prepare the microSD card

We have a separate page that describes how to [set up the microSD card]({{ "preppingthesdcard/" | prepend: site.baseurl }}) 

If you haven't done this yet, now is a good time to do it. You cannot proceed until you set up your card.
 
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Insert the card

Put the microSD card in the RPi. The slot is on the bottom of the pi. The gold-plated contacts on the card will be facing UP, and the logo on the card will be facing DOWN. Here, "UP" means "toward the top of the RPi case," and "DOWN" means "toward the bottom of the RPI case."

Also, it will only go in *one way*, so if it doesn't slot in easily, *do not force it*. Flip it around and try again if you meet resistance.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Plug in the RPi

Now it is time to plug everything in.

1. **Ethernet**. Plug one end of the ethernet cable into the RPi, and the other end into an ethernet wall jack, router, or similar.
2. **HDMI**. Plug one end of the micro-HDMI cable into the RPi (the little end) and the other end into your monitor or TV.
3. **Keyboard**. Plug your keyboard in.
4. **Mouse**. Plug your mouse in.
5. **Wifi**. Plug in the USB wifi adapter.
6. **Power**. Plug the power adapter into the switch cable (it has a cute, clicky button), and the other end of the switch cable into the RPi.

When you're done, plug the power adapter into a wall socket. *If the electrons don't flow, the Pi won't go*.
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Turn it on

Now that everything is plugged in, turn it on. Push the fun clicky switch!

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Wait

Wait. 

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Set up language and timezones

[In Wagner's video, at approximately 7:30m, you can see what the boot process should look like](https://youtu.be/7rcNjgVgc-I?t=448). 

On one hand, we would like to encourage you to select the language that is most appropriate for you. On the other, we have not tested what will happen to the data collection code if we run this with Spanish language settings (for example). Everything *should* work, but we have not *tested* it at the time of this writing. So, please choose English as the language, and you can then pick your keyboard type (it is probably a US keyboard), and your timezone.

When you're done, you should be presented with a desktop.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Test things out

Before we go further, you can test things out.

In the upper-left-hand corner is a little globe. Click on it. This is a shortcut for a web browser. In particular, *Chromium*, the open source core of Google Chrome. It will probably look familiar.

First, try visiting an exciting URL. 

[https://18f.gsa.gov/](https://18f.gsa.gov/)

That is a *very* exciting URL. Or, better yet:

[https://imls.gov/](https://imls.gov/)

If you are able to see these web pages, you are ready to proceed. [We have another page for installing the data collection software]({{ "install/" | prepend: site.baseurl }}).

If you cannot see these web pages, it means something is not working. Specifically, something involving your network. (See points of failure, below.) You now need to start trying to figure out what you may have done wrong, asking for help from colleagues, or reaching out to us. 

This is, in a way, part of the pilot. How hard is it to set these things up? How many hurdles will we run in to? What kinds of challenges will we encounter? Each is something to be understood, tackled, and solved. So, don't get frustrated; it's just part of the process. We're here for you.


# Possible points of failure

1. **Network restrictions**. It is possible your building has restrictions placed on the ethernet ports. This means that when you plug in the RPi, it will not be able to talk to the outside world. One way to figure out if this is the case is to try and do the setup in the library. If all goes well, you'll be able to open webpages using the web browser on the RPi. If you cannot, then you have two options. One is to do the setup at home. A second is to ask someone how you can get the Pi on the library's ethernet network. Note that it will need to use that network in order to take part in the pilot, so solving the underlying problem is not optional. 
1. **Missing timing on power-on**. We've noticed that, sometimes, if you turn on the RPi, and *then* turn on your monitor/TV, you might not see the RPi. We really don't know why, so we're going to make something up. It is possible that there's a synchronization that takes place while the RPi is booting, and if that sync is missed, the RPi (or the monitor/TV) don't know how to talk to each-other. (This all used to be analog. Now, it's all digital, and it really is quite complex.) We have found that if you have the TV/monitor turned on, and set to the correct input (if you have more than one DVI/HDMI input on your monitor), and *then* you turn on the RPi... then, things work pretty consistently. This may be a bunch of malarkey, but for one member of our team, it seems like this is *a thing* sometimes.