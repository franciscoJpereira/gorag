import os


CHROMA_FOLDER = "/home/franciscopereira/Documents/personalRag/chroma"
CHROMA_LOG_FILE = "chroma.log" 
CHROMA_STATE_FOLDER = "chroma_data"

# Remove directory from name
def clean_up(dirname: str):
    if not os.path.exists(dirname):
        return 
    entries = os.listdir(dirname)    
    for entry in entries:
        entry = os.path.join(dirname, entry)
        if os.path.isdir(entry):
            clean_up(entry) 
        else:
            os.remove(entry)
    os.rmdir(dirname)

if __name__ == "__main__":
    chroma_log = os.path.join(CHROMA_FOLDER, CHROMA_LOG_FILE)
    chroma_state = os.path.join(CHROMA_FOLDER, CHROMA_STATE_FOLDER) 
    clean_up(CHROMA_FOLDER)
    os.mkdir(CHROMA_FOLDER)
    os.system(f"chroma run --path {chroma_state} --log-path {chroma_log}")