# Forum


## Project overview
The forum project is the project from the Golang branch with the highest complexity. The goal was to create a forum where users could register, login, create posts with categories, see others posts, like, dislike, comment them and filter posts by categories. 

## Features

Our forum, Forum of Legends, is inspired by and named after the highly competitive online game League of Legends. It offers categorized discussions to help users connect, share insights, and grow their knowledge about the game:

**-Items:** Focuses on the game’s itemization system, providing a space for users to discuss mathematical optimizations, theorycrafting, and creative item combinations in a game with limitless strategic potential.

**-Champs:** Dedicated to the game's diverse roster of characters, this category allows users to share their experiences, strategies, and insights.

**-Meta:** Short for *Most Effective Tactics Available*, this category enables users to discuss the evolving game state, shaped by bi-weekly updates called *patches*. Here, players can explore and analyze the most current strategies and tactics to stay ahead of the competition.

**-Events:** A hub for competitive engagement, this category allows users to organize or discover tournaments and connect with teammates. Competing with stakes remains one of the fastest and most effective ways to improve.

The homepage displays the recent posts from the community disregarding categories. Different tabs helps to filter posts in a category-specific way. The unlogged user may read the posts, but cannot interact (like, dislike, comment, see others comments)

An authentification system is included in the website, with a possibility to register and add a new account to the database or create a new account. Passwords are hashed for an added layer of security. Logged users may like, dislike and comment posts, as well as create a post of their own.

The website includes a ranking system (inspired by the LoL's system to rank players) to reward active and engaging users. Posting successful content (that has a lot of likes, comments, and few dislikes) earns the user *LP* (League Points), while highly disliked posts result in LP penalties. As users accumulate LP, they rank up through the website's progression system, enhancing their credibility within the community.

A logged user can, at any time, check their created and liked posts using a special tab in the homepage, to keep track of their activity and the posts they found interesting.

## How to use it 

The long-term goal of this project might be to let the website running at all time allowing a real community to build around it. But for now, using it requires you to download the project and to run it yourself using a software such as Visual Code Studio, and you'll probably be the only user online. 

To do so, first download **Git** on your computer, then clone the git repository using ``` git clone https://zone01normandie.org/git/jpiou/forum.git ```

Make sure you have Go (version 1.20 or higher) installed on your machine. You can download Go on [their official website](https://go.dev/dl/)

Open the *forum* folder in Visual Studio Code.

Install dependencies using ```go mod tidy``` in the terminal.

Ensure you have SQLite installed on your computer, or download it on [SQLite's official website](https://www.sqlite.org/)

Please note that this program might not work on Windows. 

 Finally, use the command ``` go run . ``` in the terminal. The server is launched! Now, to access the website, just click [here](http://localhost:8080) or Ctrl + click the link in the terminal. Welcome to the Forum of Legends !

 ## How the project was done

 Starting forum was particularly challenging due to its high complexity compared to our previous projects. Building such an extensive network of interconnected programs initially felt overwhelming and discouraging. To tackle this, we chose to start small, gradually improving and refining our work to achieve higher quality over time.

 Our first step was to build the database, and make very basic html templates for the different pages we knew we would need. We had a roadmap describing everything the project would need and spent a few day brainstorming about what the website would look like.

At this stage, we developed a variety of handlers to account for the website's different functionalities, and we ended up with a working server that lacked most features. Posts could be "created" but not stored in the database, and tabs were present but couldn't filter any posts since none were registered yet. Despite its primitive state, this early version of the forum served as a stepping stone, allowing us to push through and gradually build something more substantial.

The first major roadblock we encountered was implementing the likes and dislikes system. However, as a determined team, we persevered and eventually arrived at a working solution by gaining a better understanding of primary and foreign keys in SQL and improving our data handling techniques. TThe LP system, though entirely optional and a product of our own initiative, proved also a bit challenging to implement. Still, its potential to enrich the project and enhance its overall credibility inspired us to see it through to completion.

We gradually improved the HTML templates and CSS styling as we progressed, but we hit a major hurdle with implementing comments that could be liked and disliked. After successfully creating a working version, we accidentally lost it due to a Git mishap. Undeterred, we iterated again with a barely functional version, until our most dedicated contributor rebuilt it from scratch on her own—leading to the fully working version of comments we have today!

The final major step in development was transforming a bland website into an engaging and visually appealing experience for users. We achieved this gradually, adding small bursts of code until arriving at the result we are quite proud of today.

## Usage exemple

## Future improvements

There is still plenty of room for improvement in the forum. Some of these enhancements are included as optional exercises in our school’s program, such as implementing different levels of moderation on the server, allowing certain users to edit or delete posts.

We also have the option to improve the website’s security, enable the addition of images and GIFs in posts, and much more. However, our primary goal moving forward is to enhance the website's design to make it more visually appealing. By doing so, we hope to encourage users to invest their time in the community, feel rewarded for their contributions, and create a space where individuals with shared passions can connect and engage in meaningful discussions about their favorite hobby.

## Credits

A big round of applause to the entire team for bringing this project to life in just 2-3 weeks. Given that we started our curriculum at Zone01 only two and a half months ago, tackling such a complex challenge and delivering a functional, engaging forum is an achievement we can all be proud of. Step by step, we overcame obstacles, learned along the way, and turned our efforts into something we could share.

