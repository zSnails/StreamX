import pandas as pd
from pathlib import Path
import matplotlib.pyplot as plt

# Load the CSV file
file = Path('./user-uploading-video-no-constraints.csv')

df = pd.read_csv(file)

# Convert timestamp to datetime for better readability, assuming timestamp is in UNIX format
df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')

df_grouped = df.groupby(['metric_name', 'timestamp'], as_index=False)['metric_value'].mean()

# List of unique metric names
metric_names = df_grouped['metric_name'].unique() # type: ignore

# Plot each metric separately
for metric in metric_names:
    metric_data = df_grouped[df_grouped['metric_name'] == metric]
    
    # Plot metric_value over timestamp
    plt.figure(figsize=(10, 6))
    plt.plot(metric_data['timestamp'], metric_data['metric_value'], label=metric, marker='o')

    # Add titles and labels
    plt.title(f'{metric} Over Time')
    plt.xlabel('Time')
    plt.ylabel(f'{metric} Value')
    plt.grid(True)
    plt.legend()

    plt.savefig(f'{file.stem}-{metric}_plot.svg', format='svg')

    # plt.show()

print("Plots generated for each metric.")
