import os
import sys


def main():
    if len(sys.argv) < 2:
        print("Usage: python release.py <version>")
        return
    version = sys.argv[1]
    print(f"Releasing version {version}")
    try:
        major, minor, patch = version.split('.')
        if (not major.startswith('v') and not major[1:].isdigit()) or not minor.isdigit() or not patch.isdigit():
            raise ValueError
    except ValueError:
        print("Invalid version format. Please use vX.Y.Z")
        return
    with open('version.txt', 'w') as f:
        f.write(version)
    os.system('git add version.txt')
    os.system(f'git commit -m "Bump version to {version}"')
    os.system('git push')
    os.system(f'git tag -a {version} -m "Release version {version}"')
    os.system('git push --tags')


if __name__ == '__main__':
    main()
