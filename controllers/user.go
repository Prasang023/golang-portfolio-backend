package controllers

import(
"fmt"
"context"
"log"
"net/http"
"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/bson/primitive"
"github.com/gin-gonic/gin"
)

type APIController struct{
	session *mongo.Client
}

type Member struct {
	Id string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Description string `json:"desc" bson:"desc"`
}

func NewAPIController(s *mongo.Client) *APIController{
return &APIController{s}
}

func (uc APIController) GetProjects (c *gin.Context){

	if err := uc.session.Database("portfolio-go-backend").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	coll := (uc.session).Database("portfolio-go-backend").Collection("projects"); 

	cur, err := coll.Find(context.TODO(), bson.D{})
    // cursor, err := coll.Find(context.Background(), nil)
    if err != nil {
        panic(err)
    }

	var results []Member
    for cur.Next(context.TODO()) {
        //Create a value into which the single document can be decoded
        var elem Member
        err := cur.Decode(&elem)
        if err != nil {
            log.Fatal(err)
        }

        results =append(results, elem)

    }

    if err := cur.Err(); err != nil {
        log.Fatal(err)
    }

    //Close the cursor once finished
    cur.Close(context.TODO())

    fmt.Printf("Found multiple documents: %+v\n", results)

	c.IndentedJSON(http.StatusOK, results)
	// fmt.Fprintf(w, "%s\n", uj)
}

func (uc APIController) GetProjectById (c *gin.Context){
	id := c.Param("id")

	fmt.Printf("id is: %s", id)
	objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        panic(err)
    }

    // Perform the query to find the document by ID
	coll := (uc.session).Database("portfolio-go-backend").Collection("projects"); 
	var result Member
    err = coll.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("milestone 2:", result)
	c.IndentedJSON(http.StatusOK, result)
}

// func (uc UserController) CreateUser (w http.ResponseWriter, r *http.Request, _ httprouter.Params){
// 	u := models.User{}

// 	json.NewDecoder(r.Body).Decode(&u)

// 	u.Id = bson.NewObjectId()

// 	uc.session.DB("mongo-golang").C("users").Insert(u)

// 	uj, err := json.Marshal(u)

// 	if err != nil{
// 		fmt.Println(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, "%s\n", uj)
// }


// func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params){

// 	id := p.ByName("id")

// 	if !bson.IsObjectIdHex(id){
// 		w.WriteHeader(404)
// 		return
// 	}

// 	oid := bson.ObjectIdHex(id)

// 	if err := uc.session.DB("mongo-golang").C("users").RemoveId(oid); err != nil {
// 		w.WriteHeader(404)
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, "Deleted user", oid, "\n")
// }
