package net.dev.musicchatbot;

import android.os.AsyncTask;
import android.support.design.widget.FloatingActionButton;
import android.support.v4.content.IntentCompat;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.Toast;

import org.json.*;
import com.google.gson.Gson;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;

import net.dev.musicchatbot.Adapter.CustomAdapter;
import net.dev.musicchatbot.Helper.GetOperation;
import net.dev.musicchatbot.Models.ChatModel;

public class MainActivity extends AppCompatActivity {
    boolean FirstPath = true;
    boolean Auto = false;
    int count ;
    ListView listView;
    EditText editText;
    String UUID;
    String message;
    List<ChatModel> list_chat = new ArrayList<>();
    FloatingActionButton btn_send_message;
    Button auto_button;
    final String[] MY_ARRAY = {"Hi","Human", "25", "I would like a popular song please", "lazy","thank you", "another"};

    @Override
    protected void onCreate(final Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        System.out.println("on create initiated");
        setContentView(R.layout.activity_main);

        auto_button = (Button) findViewById(R.id.auto_button);
        listView = (ListView)findViewById(R.id.list_of_message); //Gets the list containing all the messages.
        editText = (EditText)findViewById(R.id.user_message); //Gets what the user Types in the chat
        btn_send_message = (FloatingActionButton)findViewById(R.id.fab); //Floating Action Button to send

        Toast.makeText(this, "Hint type !refresh to reset the bot, say Hi!", Toast.LENGTH_LONG).show();

        //TRYING THINGS

        btn_send_message.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {

                String text = editText.getText().toString();
                if(text.contains("change mood")||text.contains("changemood")) count=3;
                if(text.contains("!refresh")) {FirstPath= true; count=0;}
                ChatModel model = new ChatModel(text,true); // user send message
                list_chat.add(model);
                new MusicBotAPI().execute(list_chat);
                count++;
                editText.setText(""); //remove user message
            }
        });

        auto_button.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                // Do something in response to button click
                String text="";
                if(count<7) text = MY_ARRAY[count];
                System.out.println(text);
                if(count>6) text="another";
                editText.setText(text);
                ChatModel model = new ChatModel(text,true); // user send message
                list_chat.add(model);
                new MusicBotAPI().execute(list_chat);
                count++;
                editText.setText(""); //remove user message
            }
        });
    }

    private class MusicBotAPI extends AsyncTask<List<ChatModel>,Void,String> {

        String stream = null;
        List<ChatModel> models;
        String text = editText.getText().toString();

        @Override
        protected String doInBackground(List<ChatModel>... params) {
            String  url = "https://floating-garden-76942.herokuapp.com/welcome";

            try {
                if(FirstPath ==true) {


                    System.out.println("Starting GET Operation");
                    String jsonString = GetOperation.getHTML(url);
                    System.out.println("finishing GET Operation");
                    FirstPath = false;
                    JSONObject jsonObject = new JSONObject(jsonString);
                    UUID = jsonObject.getString("uuid");
                    message=jsonObject.getString("message");

                }
                else{
                    System.out.println("Starting Post Operation");

                            try {

                                URL url2 = new URL("https://floating-garden-76942.herokuapp.com/chat");
                                HttpURLConnection conn = (HttpURLConnection) url2.openConnection();
                                conn.setRequestMethod("POST");
                                conn.setRequestProperty("Content-Type", "application/json;charset=UTF-8");
                                conn.setRequestProperty("Accept","application/json");
                                conn.setRequestProperty("Authorization", UUID);
                                conn.setDoOutput(true);
                                conn.setDoInput(true);

                                JSONObject jsonParam = new JSONObject();
                                jsonParam.put("uuid", UUID);
                                jsonParam.put("message", text);
                                System.out.println(UUID);
                                DataOutputStream os = new DataOutputStream(conn.getOutputStream());
                                os.writeBytes(jsonParam.toString());
                                os.flush();
                                os.close();

                                StringBuilder result = new StringBuilder();
                                Log.i("STATUS", String.valueOf(conn.getResponseCode()));
                                Log.i("MSG" , conn.getResponseMessage());
                                System.out.println("Stating return of post");
                                BufferedReader rd = new BufferedReader(new InputStreamReader(conn.getInputStream()));
                                String line;
                                while ((line = rd.readLine()) != null) {
                                    result.append(line);
                                }
                                String jsonString =result.toString();
                                JSONObject jsonObject = new JSONObject(jsonString);
                               message=jsonObject.getString("message");
                               if(message.contains("sorry"))count--;
                               if(message.contains("please"))count--;
                                System.out.println(jsonString);
                                conn.disconnect();
                            } catch (Exception e) {
                                e.printStackTrace();
                            }

                    System.out.println("Finishing Post Operation");
                }
            } catch (Exception e) {
                System.out.println("There Appears to be a problem!");
                e.printStackTrace();

            }
            models = params[0];
            stream = "Done";
            return stream;
        }

        @Override
        protected void onPostExecute(String s) {
           Gson gson = new Gson();
            ChatModel chatModel = new ChatModel(message,false); // get response from bot
            System.out.println(message);
            models.add(chatModel);
            CustomAdapter adapter = new CustomAdapter(models,getApplicationContext());
            listView.setAdapter(adapter);
        }
    }

}
