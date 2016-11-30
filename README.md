# EmergingTechnologiesProject
  We have created a single page, bloggin application in Golang, Angular JS and MongoDB.
  The website allows for the basic CRUD operations in terms of blog posts, posting comments and user registration and login authentication

##Local installation
####[Install GIT on your computer] (https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
####[Install Golang on your computer] (https://golang.org/)
####[Clone the git repository] (https://github.com/Pboyle5h/EmergingTechnologiesProject)
####Navigate to the EmergingTechnologiesProject repository folder.
####Install the gorilla toolkit with the following command:
```
go get github.com/gorilla/mux
```
####Run command:
```
go build App.go
```
####Make sure you have a working internet connection and run command:
```
go run App.go
```
####Open a web browser and go to:
Make sure port 4000 is not used by any other application.
```
localhost:4000
```
Fefer to the User Guide for further instructions.


Online:
1. Read the User Guide
2. Go to the following page : https://goproject.herokuapp.com/

User Guide:

  Registration:
  1. Click on the Register button in the top right corner.
  2. Input the data into the correct fields, in the correct format. (eg. email has to be in the "emailName"@"gmail.com" format)
  // *Limitation* The code does not currently validate whether a user with the same name already exists.
  3. Click submit
  4. You will be transported to the Login page automatically.

  Log in:
  1. Click on the login button in the top right corner.
  2. Input your login details.
  3. Click submit and await response from the database.
  ####Limitation
  '''
  Limitation The login function is currently partly working. The feedback is not correct as the user
  can be looged in correctly but it will still be displayed as invalid login. This is to do with page State
  '''

  Tabs:

    Blogs:
    1. Click on the Blogs tab in the left corner.
    2. All the blogs that are currently in the database will be displayed.
    3. Click on the blog that you are interested in reading.
    ####Limitation
    '''
    //  A comments function has been implemented in GO but no fully implemented on the HTML and Angular side
    '''
  User Page:
   1. Follow the Login Guide in order to make this tab available.
   2. Click on your username in the top right corner.
   3. A list of your personal posts will appear.

      Reviewing your blogs:
      1. Click on the blog you wish to read.
      2. Read your blog.
      // *Limitation* A comments function has been implemented in GO but no fully implemented on the HTML and Angular side

      Add a new blog:
      1. Click on the blue plus button in the top right corner.
      2. Insert a title.
      3. Insert the blog post. (Must be over 70 characters)
      4. Insert the author name.
      5. Review the blog details. The frame surrounding the text fields will turn green if the input is valid and red if it is not/ or field is empty
      6. Click the submit button at the bottom of the page.
      7. You will be redirected to the home page. Navigate back to the User Page in order to view your new post.
