package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
	"flag"
	 "math/rand"
	"google.golang.org/api/googleapi/transport"
   "google.golang.org/api/youtube/v3"
		"os"

	cors "github.com/heppu/simple-cors"
)
import _ "github.com/joho/godotenv/autoload"

var (
	 UUIDMap = make(map[string] string)
	 maxsongs int=0
	 WelcomeMessage = "Hey! What is your name?"
	 mode int8=0
	 name string
	 age int 
	 song string
	 myurl string
	 rn int
	 userrandom int64 =45
   moodArray = [26]string  {"nomood","sad","happy","relaxed","angry","excited","depressed","workout","annoyed","lazy","indifferent","fantastic","grumpy","afraid","anxious","joy",
													 "disgust","love","shame","hate","ok","bored","fine","good","awesome","tired" }
	 mood = moodArray[0]

	// sessions = {
	//   "uuid1" = Session{...},
	//   ...
	// }
	sessions = map[string]Session{}

	processor = sampleProcessor
)

const developerKey = "AIzaSyCBNMjwKAG2fS0RbcS2FOXl-i5HRK6-6c4"

type (
	// Session Holds info about a session
	Session map[string]interface{}


	// JSON Holds a JSON object
	JSON map[string]interface{}

	// Processor Alias for Process func
	Processor func(session Session, message string, uuid string) (string, error)
)

//Random Number generator
func random(min, max int) int {
    return rand.Intn(max - min) + min
}
//message is the current message from the user
func sampleProcessor(session Session, message string, uuid string) (string, error) {
	
 
	 
	// Make sure a history key is defined in the session which points to a slice of strings
	_, historyFound := session["history"]
	if !historyFound {
		session["history"] = []string{}
		UUIDMap[uuid] = "1nomood"
	}
	
	
	//Find which mode this uuid is on 
	  Current :=  UUIDMap[uuid]
	  if strings.Contains(Current, "1") { mode =1}
		if strings.Contains(Current, "2") { mode =2}
		if strings.Contains(Current, "3") { mode =3}
		if strings.Contains(Current, "4") { mode =4}
		if strings.Contains(Current, "5") { mode =5}
	
	
	
	
	// Fetch the history from session and cast it to an array of strings
	history, _ := session["history"].([]string)

		// Add the message in the parsed body to the messages in the session
	history = append(history, message)
	
	//Get The User's Name
	if mode ==1 {
			_, historyFound := session["history"]
	if !historyFound {
		session["history"] = []string{}
		sampleProcessor(session , "0",uuid)
	}

		l := len(history)
	wordsForSentence := make([]string, l)
	copy(wordsForSentence, history)

	sentence := strings.Join(wordsForSentence, ", ")

	// Save the updated history to the session
	session["history"] = history

// 			if l > 1 {
// 		wordsForSentence[l-1] = "and " + wordsForSentence[l-1]
// 	}
		  UUIDMap[uuid]="2nomood"
			return fmt.Sprintf("Hey %s! I am here to find you a music video. How old are you?", strings.ToLower(sentence)), nil
	 }
	
	//Get The User's Age
		if mode ==2 {
	
		_, historyFound := session["history"]
	if !historyFound {
		session["history"] = []string{}
		mode=0
		 UUIDMap[uuid]="2nomood"
		sampleProcessor(session , "start new session",uuid)
	}

    if _, err := strconv.Atoi(message); err != nil {
			
   return fmt.Sprintf("please write your age as a number so i can understand"), nil
					   }
			
			i, _ := strconv.ParseInt(message, 0, 64)

			if i > 100 || i < 1 {
				
				  return fmt.Sprintf("Oh come one! give me your real age so i can help you out!"), nil
			}
		l := len(history)
	wordsForSentence := make([]string, l)
	copy(wordsForSentence, history)
	session["history"] = history
			 UUIDMap[uuid]="3nomood"
			return fmt.Sprintf("Got it! Do you prefer Popular songs or would rather get something random?"), nil //Maybe we can ask what type of music he wants instead?
	}
	
	
	if mode ==3 {
		
		_, historyFound := session["history"]
	if !historyFound {
		session["history"] = []string{}
		mode=0
		sampleProcessor(session , "start new session",uuid)
	}
		
				if  strings.Contains(message, "random"){
					userrandom=50
					session["history"] = history
					 UUIDMap[uuid]="4nomood"
			return fmt.Sprintf("Random it is. what is your mood like today?"), nil //Maybe we can ask what type of music he wants instead?
					
		}
				if  strings.Contains(message, "popular"){
					userrandom=10
					session["history"] = history
					UUIDMap[uuid]="4nomood"
			return fmt.Sprintf("Popular it is. what is your mood like today?"), nil //Maybe we can ask what type of music he wants instead?
					 
		}
		session["history"] = history
		UUIDMap[uuid]="4nomood"
		return fmt.Sprintf("I didnt get that, but its ok, i am sure you will like what i pick :D \n  what is your mood like today?"), nil //Maybe we can ask what type of music he wants instead?
	 
	}
	
	//Get The User's Mood
	if mode ==4 {
		
					_, historyFound := session["history"]
	if !historyFound {
		session["history"] = []string{}
		mode=0
		sampleProcessor(session , "start new session",uuid)
	}
		mood= "nomood"
		curmood:=UUIDMap[uuid][1:len(UUIDMap[uuid])]
		fmt.Printf(curmood)
		if curmood == "nomood" {
	for j := 0; j < 26.; j++ {
		if  strings.Contains(message, moodArray[j]){
			
			mood= moodArray[j]
			curmood=mood
			 UUIDMap[uuid]="5"+mood
			break
		}
	}
		}
	l := len(history)
	wordsForSentence := make([]string, l)
	copy(wordsForSentence, history)


	// Save the updated history to the session
	session["history"] = history
		if curmood!="nomood"{
			//YOUTUBE SEARCH FUNCTION:
		
			//var query  = flag.String("query", mood + " song", "Search term")
		 flag.Parse()

        client := &http.Client{
                Transport: &transport.APIKey{Key: "AIzaSyCBNMjwKAG2fS0RbcS2FOXl-i5HRK6-6c4"},
        }

        service, err := youtube.New(client)
        if err != nil {
                log.Fatalf("Error creating new YouTube client: %v", err)
        }

        // Make the API call to YouTube.
        call := service.Search.List("id,snippet").
                Q( curmood + " song").
                MaxResults(userrandom)
        response, err := call.Do()
        if err != nil {
                log.Fatalf("Error making search API call: %v", err)
        }

        // Group video, channel, and playlist results in separate lists.
        videos := make(map[string]string)
        channels := make(map[string]string)
        playlists := make(map[string]string)

        // Iterate through each item and add it to the correct list.
        for _, item := range response.Items {
                switch item.Id.Kind {
                case "youtube#video":
									maxsongs++
									videos[item.Id.VideoId] = item.Snippet.Title
                case "youtube#channel":
                        channels[item.Id.ChannelId] = item.Snippet.Title
                case "youtube#playlist":
                        playlists[item.Id.PlaylistId] = item.Snippet.Title
                }
        }
			  rn = random(1, maxsongs)
        printIDs("Videos", videos)
       // printIDs("Channels", channels)
      //  printIDs("Playlists", playlists)
			//return fmt.Sprintf(playlists), nil 
			//return fmt.Sprintf("All done here!"), nil
			//u, err := url.Parse(myurl)
			UUIDMap[uuid]="5"+curmood
			return fmt.Sprintf("Gotcha! you are feeling %s! how is this song? \n"+song +"\n"+ myurl , strings.ToLower(curmood)), nil
			
		} 
		 UUIDMap[uuid]="4nomood"
		return fmt.Sprintf("I am sorry i don't understand your mood :( How about we try again?"), nil
	}
	if mode == 5 {
		
		if  strings.Contains(message, "another"){
					curmood:=UUIDMap[uuid][1:len(UUIDMap[uuid])]
			 UUIDMap[uuid]="4"+curmood
			maxsongs=0
			return sampleProcessor( session,  "another",uuid)
					//	return fmt.Sprintf("Right away! "), nil
		} else if  strings.Contains(message, "mood"){
	
					maxsongs=0
		mood = "nomood"
			 UUIDMap[uuid]="4nomood"
		//	curmood="nomood"
			session["history"] = "hi"
							return fmt.Sprintf("I hope your mood has improved! how are you feeling now? "), nil
		} 
		
		
				l := len(history)
	wordsForSentence := make([]string, l)
	copy(wordsForSentence, history)
	session["history"] = history
		curmood:=UUIDMap[uuid][1:len(UUIDMap[uuid])]
		UUIDMap[uuid]="4"+curmood
			return fmt.Sprintf("If you want a different song please type \"another\" or if your mood has changed type \"change mood\"  "), nil
		
	}
	
	
	
	return fmt.Sprintf("All done here!"), nil
}


// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
        fmt.Printf("%v:\n", sectionName)
        for id, title := range matches {
					rn--
					if rn ==0 {
						song = title
						myurl = "https://www.youtube.com/watch?v="+id
					}
                fmt.Printf("[%v] %v\n", "https://www.youtube.com/watch?v="+id, title)
        }
        fmt.Printf("\n\n")
}

// withLog Wraps HandlerFuncs to log requests to Stdout
func withLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := httptest.NewRecorder()
		fn(c, r)
		log.Printf("[%d] %-4s %s\n", c.Code, r.Method, r.URL.Path)

		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		c.Body.WriteTo(w)
	}
}

// writeJSON Writes the JSON equivilant for data into ResponseWriter w
func writeJSON(w http.ResponseWriter, data JSON) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ProcessFunc Sets the processor of the chatbot
func ProcessFunc(p Processor) {
	processor = p
}

// handleWelcome Handles /welcome and responds with a welcome message and a generated UUID
func handleWelcome(w http.ResponseWriter, r *http.Request) {
	// Generate a UUID.
	hasher := md5.New()
	hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	uuid := hex.EncodeToString(hasher.Sum(nil))

	// Create a session for this UUID
	sessions[uuid] = Session{}

	// Write a JSON containg the welcome message and the generated UUID
	writeJSON(w, JSON{
		"uuid":    uuid,
		"message": WelcomeMessage,
	})
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	// Make sure only POST requests are handled
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Make sure a UUID exists in the Authorization header
	uuid := r.Header.Get("Authorization")
	if uuid == "" {
		http.Error(w, "Missing or empty Authorization header.", http.StatusUnauthorized)
		return
	}

	// Make sure a session exists for the extracted UUID
	session, sessionFound := sessions[uuid]
	if !sessionFound {
		http.Error(w, fmt.Sprintf("No session found for: %v.", uuid), http.StatusUnauthorized)
		return
	}

	// Parse the JSON string in the body of the request
	data := JSON{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Couldn't decode JSON: %v.", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Make sure a message key is defined in the body of the request
	_, messageFound := data["message"]
	if !messageFound {
		http.Error(w, "Missing message key in body.", http.StatusBadRequest)
		return
	}

	// Process the received message
	message, err := processor(session, data["message"].(string),uuid)
	if err != nil {
		http.Error(w, err.Error(), 422 /* http.StatusUnprocessableEntity */)
		return
	}

	// Write a JSON containg the processed response
	writeJSON(w, JSON{
		"message": message,
	})
}

// handle Handles /
func handle(w http.ResponseWriter, r *http.Request) {
	body :=
		"<!DOCTYPE html><html><head><title>Chatbot</title></head><body><pre style=\"font-family: monospace;\">\n" +
			"Available Routes:\n\n" +
			"  GET  /welcome -> handleWelcome\n" +
			"  POST /chat    -> handleChat\n" +
			"  GET  /        -> handle        (current)\n" +
			"</pre></body></html>"
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintln(w, body)
}

// Engage Gives control to the chatbot
func Engage(addr string) error {
	// HandleFuncs
	mux := http.NewServeMux()
	mux.HandleFunc("/welcome", withLog(handleWelcome))
	mux.HandleFunc("/chat", withLog(handleChat))
	mux.HandleFunc("/", withLog(handle))

	// Start the server
	return http.ListenAndServe(addr, cors.CORS(mux))
}




// Autoload environment variables in .env


// func chatbotProcess(session chatbot.Session, message string) (string, error) {
// 	if strings.EqualFold(message, "chatbot") {
// 		return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
// 	}

// 	var questionMarksCount int
// 	// Try fetching the count of question marks
// 	count, found := session["questionMarksCount"]
// 	// If a count is saved in the session
// 	if found {
// 		// Cast it into an int (since sessions values are generic)
// 		questionMarksCount = count.(int)
// 	} else {
// 		// Otherwise, initialize the count to 1
// 		questionMarksCount = 1
// 	}

// 	// Build the question marks string according to the question marks count
// 	var questionMarks string
// 	for i := 1; i <= questionMarksCount; i++ {
// 		questionMarks += "?"
// 	}

// 	// Save the updated question marks count to the session
// 	session["questionMarksCount"] = questionMarksCount + 1

// 	// Return the response with an extra question mark
// 	return fmt.Sprintf("Hello <b>%s</b>, my name is chatbot. What was yours again%s", message, questionMarks), nil
// }

func main() {
	// Uncomment the following lines to customize the chatbot
	// chatbot.WelcomeMessage = "What's your name?"
	// chatbot.ProcessFunc(chatbotProcess)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "3000"
	}

	// Start the server
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatalln(Engage(":" + port))
}
