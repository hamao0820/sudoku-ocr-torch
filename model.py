from torch import nn
from torch.nn import functional as F


class OCRModel(nn.Module):
    def __init__(self):
        super(OCRModel, self).__init__()
        # input 3x64x64
        self.conv1 = nn.Conv2d(3, 32, 3, padding=1)
        # 32x64x64
        self.conv2 = nn.Conv2d(32, 64, 3, padding=1)
        # 64x64x64
        self.pool = nn.MaxPool2d(2, 2)
        # 64x32x32
        self.fc1 = nn.Linear(64 * 16 * 16, 128)
        self.fc2 = nn.Linear(128, 10)

        self.dropout = nn.Dropout(0.5)
        self.flatten = nn.Flatten()

    def forward(self, x):
        x = self.pool(F.relu(self.conv1(x)))
        x = self.pool(F.relu(self.conv2(x)))
        x = self.flatten(x)
        x = F.relu(self.fc1(x))
        x = self.dropout(x)
        x = self.fc2(x)
        x = F.log_softmax(x, dim=1)
        return x
