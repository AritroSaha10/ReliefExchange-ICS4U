<div align="center">
<img width="300" src="https://github.com/AritroSaha10/ReliefExchange-ICS4U/assets/76495965/ed1cc35c-ec39-4066-b8a9-e0a5d53821d2" />
</div>

# ReliefExchange: ICS4U Final Project
ReliefExchange is a platform where **Generosity meets Community**. This website built by students at John Fraser Secondary School is meant for newcomers, immigrants, or any hard working individuals to find the essentials they need. In a society where individuals are struggling to be able to simply get the everyday items they need, and where waste is abundant, we built a platform where users can donate and trade for absolutely free! On the Relief Exchange platform, users sign up to create an account that allows them to either post an item they don’t need anymore, or browse through almost any category to find something for them. With how easy it is to simply throw something away, traditional donation centres are receving less and less items each year. By having an easy to use, online platform, we hope promote accessibility so users are more motivated to donate their items, instead of aiding to our ecological impact. To sign up now, visit https://reliefexchange.aritrosaha.ca/ 

# Features 
* User friendly interface that does not require the user to be well experienced with apps & Websites
* Ability to translate the website in over 6 different languages for immigrants and non-english speaking individuals
* Search to find Donation offers using filters including tags/categories, date uploaded, and more
* Upload a donation item with descriptions, images, location, and tags 
* Profile that tracks your donations, and users registration and sign-in date 

# Built With 
- TypeScript
    - React
    - Next.js
- Golang
    - Gin
- Firebase
    - Firestore
    - Authentication
    - Storage
- Vercel (Frontend)
- Oracle Cloud Interface (Backend)
    - 1x VM w/ 4 Ampere ARM cores
    - Running as Docker container
    - NGINX (reverse proxy for container)
- GitHub Actions (CI/CD)
    - Both using GitHub-hosted servers and self-runners on OCI
- GitHub Container Registry (CI/CD)
- Testing Frameworks
    - React Testing Library
    - Jest
    - Cypress
    - `go-test`
    - `mockfs`
- Sentry (error logging for full stack)

# Development Instructions
 Make sure to change your working directory to the appropriate folder before running the commands.

 ## Frontend
 - Install dependencies using `npm install`
 - Copy `.env.template` into `.env.local`, and populate the environment variables with your values
 - Run the development server using `npm run dev`
 - Build the production code using `npm run build`
 - Run the production server locally using `npm run start`
 - Run tests using `npm run test` and `npm run cypress`

 ## Backend
 - Install dependencies using `go mod download && go mod verify`
 - Copy `.env.template` into `.env`, and populate the environment variables with your values
 - Run the development server using `go run .`
 - Build the production binary using `go build .`
 - Run the production binary as you would any other executable binary file
 - Run tests using `go test`

 # Deployment Instructions

 ## Frontend
 You have a few choices to deploy the frontend code.
 - Deploy to Vercel (what we do) or Netlify (not tested)
 - Deploy to another server (learn more about this [here](https://nextjs.org/docs/pages/building-your-application/deploying))
     - Please note that Static HTML Export is not possible for this project as we use Incremental Static Regeneration (ISR).

 ## Backend
 A Dockerfile and example Docker Compose file is available in the root directory of the backend folder. You can deploy this on any AMD64 or ARM64 machine. The built Docker image is available at [ghcr.io/aritrosaha10/relief-exchange-backend](https://github.com/aritrosaha10/ReliefExchange-ICS4U/pkgs/container/relief-exchange-backend). 
 
Our deployment runs this image as a singular Docker container on Oracle Cloud Infrastruction (OCI) behind a NGINX reverse proxy. It is highly recommended to use a reverse proxy such as Traefik or NGINX.

 ## Firebase
 - Initialize a new Firebase project
 - Add your own configuration data to the appropiate frontend and backend
 - Set up:
     - Firestore
         - Copy-paste the security rules found [here](firebase/rules/firestore.cel)
         - Create three collections:
             - `config`
                 - This should have one document called `bans`
                     - This should have an array with the key of `users`
             - `donations`
             - `users`
     - Authentication
         - Add Google as a sign-in provider
     - Storage
         - Copy-paste the rules found [here](firebase/rules/storage.cel)
         - Create a folder called `donations`

# User Testing Results
Our user testing results can be seen in [this document](USER_TESTING.md).

# How to Navigate Github
- In the navbar options below the repository title, the user can view the repository code, issues (assignments), and pull requests (new code requested to be added to the main branch of the project)
- **<>Code tab:** In the code tab, the user can view the files and folders of the project that contain the project code by clicking on the folder or file icons. On the same page, the user can also view 
- Branches: On the left side of the <>Code tab, you will see a dropdown with the label of the branch name, in our case it should say "master". By clicking on this dropdown, you can view the branches of the repository and switch between them. Each branch represents a different version or part of the project that is separately maintained.

-Commits: On the right side of the screen, under the green "Code" button, you will see an option labeled "commits". Click on it to view all commits. Each commit represents a saved change to the code, and clicking on it will show you what changes were made, who made them, and when.

-**Issues tab:**  Contains the work breakdown of who should work on what assignment. The issue title is listed on the left and the assignee's profile picture is on the right. Issues also allow you to view reported bugs or feature requests. Issues are a way for contributors to communicate about the project. If you click on an issue, you'll see a discussion about it and any associated commit or pull request.

-**Pull Requests tab:** To view proposed changes to the code (pull requests), click on the "Pull requests" tab in the navbar. Here, you can view open and closed pull requests, review the proposed changes, and participate in discussions about those changes. If you have write access to the repository, you can also merge pull requests here.

-**Actions:** The Actions tab is where you can view the history of all the tasks that have been run in the repository, often for testing or deploying the code.

-**Projects:** In the Projects tab, you can create and view project boards. These boards can be used to manage tasks, plan future work, or track progress on ongoing work.

-**Wiki:** The Wiki tab can be used to host documentation for the project. It’s a good place to look for more detailed information about how to use the project or contribute to it.

-**Security:** The Security tab provides an overview of the security policies of the project. It's also where you can view or report security vulnerabilities.

-**Insights:** Under the Insights tab, you will find statistics and analytical data related to the project. This includes things like a graph of commit activity, a list of contributors, and a dependency graph.

-**Settings:** The last tab, Settings, is only visible to repository administrators. This is where various repository settings can be changed, such as the repository's name, description, and visibility; the branch used for GitHub Pages; and more.


# Authors
This project was worked on by Aritro Saha, Joshua Chou, and Taha Khan as a part of our ICS4U0 (Computer Science) final project. 
