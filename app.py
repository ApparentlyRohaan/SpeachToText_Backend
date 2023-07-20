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

# Specify the path to your audio file
audio_file = "file.wav"
# audio_file = "uploads/recording.webm"

# Convert speech to text
transcript = speech_to_text(audio_file)

# Print the transcript
print(transcript)
