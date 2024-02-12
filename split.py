from glob import glob
import os
import random
import shutil


train_ratio = 0.8
valid_ratio = 0.1

# データセットのファイルパスを取得
dirs = glob("data/cells/**")
for d in dirs:
    base = os.path.basename(d)
    files = glob(f"{d}/*")
    files_indices = list(range(len(files)))
    train_size = int(len(files) * train_ratio)
    valid_size = int(len(files) * valid_ratio)
    test_size = len(files) - train_size - valid_size
    random.shuffle(files_indices)
    train_indices = files_indices[:train_size]
    valid_indices = files_indices[train_size : train_size + valid_size]
    test_indices = files_indices[train_size + valid_size :]
    train_files = [files[i] for i in train_indices]
    valid_files = [files[i] for i in valid_indices]
    test_files = [files[i] for i in test_indices]
    os.makedirs(f"data/train/train/{base}", exist_ok=True)
    os.makedirs(f"data/train/valid/{base}", exist_ok=True)
    os.makedirs(f"data/train/test/{base}", exist_ok=True)
    f_name = 0
    for f in train_files:
        shutil.copy(f, f"data/train/train/{base}/{f_name}.png")
        f_name += 1
    f_name = 0
    for f in valid_files:
        shutil.copy(f, f"data/train/valid/{base}/{f_name}.png")
        f_name += 1
    f_name = 0
    for f in test_files:
        shutil.copy(f, f"data/train/test/{base}/{f_name}.png")
        f_name += 1
