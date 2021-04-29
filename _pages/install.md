---
title: Install and Config
layout: page
sidenav: false
---
{% assign counter = 0 %}

# {{page.title}}

Here we go! You have an RPi, it is set up, and it can talk to the internet! This is very, *very* exciting.

##  Step {% assign counter = counter | plus: 1  %}{{counter}}: Visit this page on the RPi

You really need this page open on the Raspberry Pi. We've given you a short URL you can type in:

[https://go.usa.gov/xHQYj](https://go.usa.gov/xHQYj)

That's `https://`, then `go.usa.gov`, a slash, and then `x` `H` `Q` `Y` and `j`. 

Because the next step involves some copy-paste work, you need to have this page open on the Pi. 
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Open the terminal

In the upper-left corner of the screen is a little black box with a blue bar at the top. This icon is intended to represent a command terminal. It should open up a window that looks like this:

{% asset "pi-command-prompt.png" alt="location of pi command prompt icon" %}

Go ahead and grab the bottom-right corner and make the window a bit bigger before continuing.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Copy the line below

```
bash <(curl -s https://raw.githubusercontent.com/jadudm/imls-client-pi-playbook/main/bootstrap.sh)
```

Highlight that line, right-click, and say **Copy**. 

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Paste the line into the terminal

Over in the terminal, right-click, and select **Paste**. The command should now be in the terminal window.

Press Enter. Or, Return. Or whatever that key is called on your keyboard.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter your FCFS_Seq_Id

You are approaching the final stretch.

You are setting up your device so that it will live in a particular library. Or, perhaps, a bookmobile. Either way, there should be an FSFS_Seq_Id for your location.

You can [look this up your FCFS Seq Id for your building on the IMLS website](https://www.imls.gov/search-compare/). For example, the Berea branch in Madison County, KY is **KY0069-003**.

Type this in when prompted. You will need to use caps for the state abbreviation (e.g. you cannot say `ky`, but you can say `KY`).

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter a tag

This is a free-form text entry that will accept up to 255 characters. It *will* become part of the dataset, so please do not put any personally identifiable information into this field. Email addresses are considered PII.

We intend this field either as a "hardware tag" (so you might enter `RPi 001`), or it could be a reminder where you put it. In the latter case, you might enter `reference desk`, or `networking closet`, or (if it is in the bookmobile)... er, perhaps `glovebox`? Point being, this is mostly for you. This will become part of the dataset, so something simple but descriptive is best.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter your api.data.gov key

This is a tedious part.

First, note: **you should never, ever give anyone your API key**. You should not paste it into websites. You should not mail it to your friends. (They keys are free... they can go get their own.) But, we're going to ask you to violate that rule right now.

On your other computer (not the Pi), you'll want to find the email where you received your API key from api.data.gov. Then, you'll want to visit this page. (Yes, this one, right here.) Then, you'll want to paste your key into the text box below:

<script type="text/javascript" src="{{ '/js/wordlist.js' | prepend: site.baseurl }}"></script>

<div class="grid-container">
  <div class="grid-row">
    <div class="tablet:grid-col">
        <label class="usa-label" for="api-key">Your API Key</label>
        <input class="usa-input" id="api-key" name="api-key" type="text" placeholder="library">
    </div>
    <div class="grid-col-fill" style="padding-left: 2em;">
        <table class="usa-table usa-table--borderless usa-table--striped">
            <thead>
                <tr>
                    <th scope="col">Words</th>
                    <th scope="col">Decoded</th>
                </tr>
            </thead>
            <tbody id="tablebody">
            </tbody>
        </table>
    </div>
</div>
</div>

Your key will not be sent anywhere. The transformation that takes place happens *entirely within this webpage*. That is the reason it is safe for you to paste your API key here (and nowhere else!). 

Your API key has been transformed into a list of "word pairs."

You'll read each word-pair (perhaps off your phone, or a laptop off to the side) and type them into the Raspberry Pi one pair at a time. The Pi will make sure that you're typing "correct" word pairs. (There's a limited set; if you make a typo, we'll ask you to try again.) You should end up entering 14 word pairs.

As you type your word pairs, the Raspberry Pi will "translate" those words back into pieces of your API key. 

When you are done, you can type `DONE`. Or, `done`. 

This will now save the API key on the Raspberry Pi.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Installation continues

At this point, the Pi is going to do a whole bunch of things on its own. This process is *slow*. But, you don't have to do anything more; the rest is fully automated.

What is happening? This is where we take the friendly Raspberry Pi that you were sent, and we turn it into a secure device. Amongst other things, we are:

1. Installing software to do the WIFISESS data collection. 
2. Making sure the RPi gets security and software updates. 
3. Disabling the ability for anyone to connect to the RPi via the network.
4. Disabling the ability for anyone to log into the RPi with a keyboard and mouse.
5. Removing the graphical interface.

And a few more things. When we're done, you really have no idea if it worked, save for the fact that it will reboot, and then quietly do its job. 

##  Step {% assign counter = counter | plus: 1  %}{{counter}}: Is it working?

FIXME

*THIS WOULD BE A GREAT PLACE TO EMBED A JAVASCRIPT DEMO*



{% asset token.js %}
