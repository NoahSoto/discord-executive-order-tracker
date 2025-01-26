
from sumy.parsers.plaintext import PlaintextParser
from sumy.nlp.tokenizers import Tokenizer
from sumy.summarizers.lsa import LsaSummarizer

# Input text to be summarized
f = open("order.txt" ,"r")
input_text = f.read()
# Parse the input text
parser = PlaintextParser.from_string(input_text, Tokenizer("english"))

# Create an LSA summarizer
summarizer = LsaSummarizer()

# Generate the summary
summary = summarizer(parser.document, sentences_count=3)  # You can adjust the number of sentences in the summary

# Output the summary
with open("summary.txt", "w") as file:
    pass  # This clears the file by overwriting it with an empty file

with open("summary.txt", "a") as file:
    for sentence in summary:
        file.write(str(sentence))
