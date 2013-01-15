Music playback philosophy
-------------------------

Don't mind me, i'm just thinking loudly.

Music should have, additionally to being ordered in Artist/Album/Genre, attributes like _mood_, _tempo_,
_syntheticness_ or _aggressivity_.

The concept of tags, which I quite like, is an approach to this. By having a ginormous amount of sub-genres, sub-sub-genres and sub³genres,
giving a song tags can be time-consuming and mostly inaccurate. _Last.fm_ solved this by having each user tag his/her songs
and display the most-used tags on a song first. For this to work, the song has to be listened to by many users.
For less well-known music, there is never enough tag-data for a song.

How to solve this?

### With Math
You heard right. You could describe a song as a set of properties (see above) and their strength.
But: assigning these values (for me) often requires more thinking time than just assigning genres (or subgenres).
Genres are basically pre-defined sets of values for the abovementioned properties, and because they are easier to deal with
(in your brain), let's use a set of genres and a value for how much this song represents a genre.

Mathematically, a song is just a point in the n-dimensional room of properties. Viewing it as a room where the genres
are the dimensions (so just a different [basis](http://en.wikipedia.org/wiki/Basis_%28linear_algebra%29)), we can assign it
coordinates in that space. If we have the mappings from our genres to their properties, using a little geometry, we can
derive a songs properties by transforming the (genre,value) tuples to the (property,value) base.

The number of dimensions is rather infinite, let's limit it to the number of genres a user listens to.

### What now?
So now that we have every song assigned a vector in our genre-base, we can calculate the songs properties and also know,
which tags (= genres) other users would assign to such a song.

Using the data we gain about the song from the user, we can construct and correct the genres position in the property-room,
assigning each genre a vector of properties.
Going even further, since genres are often broad definitions, we can instead assign to a genre for each property a mean value and a standard derivation,
which, when combined, gives us a definition of the variations (and nonvariations) of a genre.

###We need to go deeper
Additionally, a user might react differently to a given song based on his current mood, situation or the time,
giving it different ratings. So, when assigning values to a song, we need to check how the user rated the other songs
he listened to in the last few hours, comparing offsets and trends in specific properties.

If we ask the user about his mood, we can even build a profile and, after there is enough data, tell him what genres he might
like in his current mood.

---
to be continued.


_Conclusion:_ iTunes just isn't good enough.
