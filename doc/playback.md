Music organizing philosophy
---------------------------

Some general thoughts about how to store music metadata.

Music is most fun when it's surprising, but meets your taste. Music listening (marketing) systems like Last.fm
are based on that principle. Underlying suitable algorithms thus need a formalization of aesthetical similarity
and a concept of what kind of music it is.

Music in a database should therefore, additionally to being ordered in Artist/Album/Genre, have attributes like _mood_, _tempo_,
_syntheticness_ or _aggressiveness_.  
The concept of tags, which I quite like, is an approach to this. Tags are boolean only, so they can't always be precise.
A song can have a tag or -- well, -- not have a tag.  
By having a ginormous amount of sub-genres, sub-sub-genres and sub³genres, one can still add some granularity.
However, giving a song such tags can be time-consuming and mostly inaccurate.  
_Last.fm_ solved this by having each user tag his/her songs and display the most-used tags on a song first.
For this to work, the song has to be listened to by many users.
For less well-known music, there is never enough tag data for a song.

How to solve this?

### With Math
Woo, math. You could describe a song with a set of properties (see above)
each having a value specifying to what extent that property applies.  
Assigning these values (for me) often requires more thinking time than just assigning genres (or subgenres).
Genres are basically predefined sets of values for the abovementioned properties, and because they are easier to deal with
(in your brain), let's use a set of genres and a value for how much this song represents a genre.

Mathematically, a song is just a point in the n-dimensional space of properties. Viewing it as a space where the genres
are the dimensions (so just a different [basis](http://en.wikipedia.org/wiki/Basis_%28linear_algebra%29), linear algebra whee),
we can assign it coordinates in that space.  
If we have the mappings from our genres to their properties, using a little geometry, we can
derive properties from a song by transforming the (genre,value) tuples to the (property,value) basis.
This basis change happens through a transformation matrix, let's call it _M<sub>g&rarr;p</sub>_.

The number of dimensions is rather infinite, let's limit it to the number of genres a user listens to.

### What now?
So now that we have every song assigned a vector in our genre-basis, we can calculate the songs properties and also know,
which tags (= genres) other users would assign to such a song.

Using the data we gain about the song from the user, we can construct and correct the genres position in the property-room,
assigning each genre a vector of properties. From that, we can also construct _M<sub>g&rarr;p</sub>_.

Going even further, since genres are often broad definitions, we can instead assign to a genre for each property a mean value and a standard derivation,
which, when combined, gives us a definition of the variations (and fixed properties) of a genre.

###We need to go deeper
Additionally, a user might react differently to a given song based on his current mood, situation or the time,
giving it different ratings. So, when assigning values to a song, we need to compensate for how the user rated the other songs
he listened to in the last few hours, comparing offsets and trends in specific properties.

If we ask the user about his mood, we can even build a profile and, when there is enough data, tell him what genres he might
like in his current mood.

---
Next time: How do we capture all that data and use it to play the user music to his liking?


_Conclusion:_ iTunes just isn't good enough.
