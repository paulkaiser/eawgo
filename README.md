# Empires At War

## The Game
### Territories
Territories are the basis for dividing the known world into controllable areas. 

*	A Territory is known by its name. 
*	A Territory has exactly one Player who owns. 
*	A Territory has a list of one or more attack vectors. An attack vector is a Territory that can be attacked from the given Territory.

### Players
Players represent the Emporers who are battling for control of Territories.

*	A Player is known by order and name.
*	A player owns one or more Territories.
*	A player without Territories cannot play.

## Running the Game
### Starting the Game
### Rounds
### Player Turns
### Ending the Game

## Technical Notes
### Communication Structure
Goal: use channels to abstract communication between main game thread, game play controller and users.

Each Player has three channels: 

*	one for asynchronous output (game-to-player), i.e., output for informational or instructional purposes
*	one for asynchronous input (player-to-game), i.e., input that does not have to be solicited to be used
*	one for request/response interactions, i.e., output prompts to the user with synchronous input expected

Start with something simpler, but define channels at a functional level. They should be used bi-directionally without problems with unintended inputs that would otherwise require filtering.

May need to have main game thread control whether a new round is started (requires consensus of all players) and then delegate control of the entire round to the game play controller.

## Backlog
### Core game 

#### Add territory state to UX sequence
After each attack, write the outcome to the console and show the current state of territory ownership before prompting "Do you want to attack?"

#### Improve randomness of dice rolls
Currently, the dice package seeds the random number generator using a value from the environment or a simple default value. The apparent randomness can be improved by using a variation of clock values to form the seed.

#### Starting the game
Select the number of players and provide names for each.

#### Players select territories
Change the assignment of territories to players. Once the order of play is established, let each player select one territory at a time, in the order of play, until all territories are claimed. All territories must be claimed even if some players have one more territories than others.

#### Ending the game
The length of the game can be defined at the beginning 

* as a specific number of rounds
* with a time limit
* unlimited

If the game is to run a specific number of rounds, the game is over after each player has completed that number of rounds. Players may still withdraw from the game per the withdrawal rules. This does not affect the number of rounds played.

If the game is to end after a specific amount of time has elapsed, the round in play when the time expires will be completed. Players may still withdraw from the game per the withdrawal rules.

If an unlimited game is selected, play will continue until the first player refuses to start a new round and does not chose to withdraw.

#### The Winner
The player with the most territories at the end of the game wins.
Broadcast the winner when the game is over.

#### The Leader
The player with the most territories at the end of a round is the leader.
Broadcast the leader when the game is over.

#### Evaluate managing players in a ring
Consider using a ring as the data structure that represents the order of play instead of an array. The goal should be to simplify managing the order of play.


#### Improve reaction to invalid user input
Given	the user is presented with a list of choices
and 	the system prompted for a selection
and		the user enters a choice not in the list
then	the system will report the error
and		the system will present the list of choices
and		the system will prompt for a selection


#### Put a log level wrapper around the log package
That just about says it all.

Use lumber - https://github.com/jcelliott/lumber


#### Add support for backing out of an attack
Given	the current player has chosen to attack
but		has not selected both attacking and defending territories
then	the current player may call off the attack.

#### Let a player suspend play
TBD

#### Let players suspend game
TBD

#### Let a player withdraw from the game
TBD

### Game server
#### Core game encapsulation
Build a socket-based server that encapsulates an instance of the core game. The goal is to let the remote client application initiate games, enlist players and start a game. The core game must be able to run without the game server using whatever user interaction facilitities inherent in the core. 

The thinking is that channel-based player interaction is the key abstraction to making this work.


### Client application
#### Support simple remote player interactions
Building on top of the channel-based conversational paths with users, build a simple application that connects to a game server

#### Display of current territory state
Once the game has been updated to broadcast attack outcomes and territory updates after each attack, the client application should render it simply. The simple visual can be a list of

* territory name
* owner
* list of attack vectors
	* each territory that can be attacked
	* its owner
