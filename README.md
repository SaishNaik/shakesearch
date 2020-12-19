# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://shake-search-great.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is just a rough prototype. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the search backend. Think about the problem from the user's perspective
and prioritize your changes according to what you think is most useful.

To submit your solution, fork this repository and send us a link to your fork
after pushing your changes. The project includes a Heroku Procfile and, in its
current state, can be deployed easily on Heroku's free tier.

If you are stronger on the front-end, complete the react-prompt.md in this
folder.

## Changes  

1) Case insensitive search.
2) UI improved with highlighting of the text selected and better presentation.
3) Meaningful slices of text shown for the word searched rather than the fixed 250 length before.
So there wont be any sliced data.
4) Exact match issue solved by searching substrings from front and back. Eg. if "hamlet" is searched, but is not available, it will return results for "amlet" "mlet","let". It will also search for "hamle","haml","ham".
5) Incase of unmatched string, it will search till 3 characters are remaining in the query. eg for "hamlet", it will search till  "ham".
6) Searching for exact matches can be a costly operation and hence go channels are used for parallelism. 
7) No Duplicate Matched text. Due to exact matching implemented, This issue is likely to occur. Fixed.
