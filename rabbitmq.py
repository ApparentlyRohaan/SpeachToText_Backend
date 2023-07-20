import pika
import speech_recognition as sr


def speech_to_text(audio_file):
    # Create a recognizer instance
    recognizer = sr.Recognizer()

    # Load the audio file
    with sr.AudioFile(audio_file) as source:
        audio = recognizer.record(source)

    try:
        # Perform speech recognition
        text = recognizer.recognize_google(audio)
        return text
    except sr.UnknownValueError:
        print("Speech recognition could not understand audio")
    except sr.RequestError as e:
        print("Could not request results from Google Speech Recognition service; {0}".format(e))

    return ""


def callback(ch, method, properties, body):
    message = body.decode()
    print(f"Received: {message}")

    # Prepare the response message
    
    # Specify the path to your audio file
    audio_file = "file.wav"
    # audio_file = "uploads/recording.webm"

    # Convert speech to text
    transcript = speech_to_text(audio_file)

    # Print the transcript
    print(transcript)

    response = f"Python received: {transcript}"


    # channel.basic_publish(exchange='',
    #                       routing_key=method.routing_key,
    #                       body=response,
    #                       properties=pika.BasicProperties(
    #                           delivery_mode=2,  # Make the message persistent
    #                       ))



    # print(f"Sent response: {response}")

connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
channel = connection.channel()

queue_name = 'hello'
channel.queue_declare(queue=queue_name)

channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)

print('Waiting for messages. To exit press CTRL+C')
channel.start_consuming()
