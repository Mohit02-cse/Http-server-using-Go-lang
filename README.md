Shopping List API with MySQL Integration
<br>
This project is a RESTful API developed in Go (Golang) that allows users to manage shopping lists for different customers. The application provides functionality to 
create, retrieve, and delete shopping items, while persisting data in a MySQL database.

Key features of the project include:

    Dynamic Database Management:
        The API ensures that a database named Shopping_List is created if it does not already exist.
        Separate tables are dynamically created for each customer to organize their shopping lists efficiently.

    CRUD Operations:
        Create: Users can add items to a customer's shopping list.
        Read: Users can view all items in a customer's shopping list.
        Delete: Users can remove specific items from a shopping list using their unique IDs.

    Robust Error Handling:
        The application handles invalid JSON payloads, missing or incorrect data, and SQL errors gracefully, providing clear feedback to the client.

    Integration with MySQL:
        Utilizes Go's database/sql package to interact with the MySQL database.
        Ensures schema creation and data integrity during runtime.

    Lightweight and Scalable:
        Built using the gorilla/mux package for routing, making the API modular and easy to extend.

Technologies Used

    Programming Language: Go (Golang)
    Database: MySQL
    Libraries/Frameworks:
        gorilla/mux for routing
        github.com/google/uuid for unique ID generation
    Development Tools: Git, Postman (for testing)

This project demonstrates practical use of Go for backend development, efficient database integration, and dynamic table management for scalable and user-specific data handling.
