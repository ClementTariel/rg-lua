import sys

if (len(sys.argv) != 2 or not sys.argv[1].isdigit()):
    print("usage:\npython3 generate_grid.py <RADIUS>")
    exit()

output_file = open("grid.c", "w")

radius = int(sys.argv[1])
grid_size = 2*radius + 3

center = (grid_size - 1) / 2
min2 = (radius - 0.5) * (radius - 0.5)
max2 = (radius + 0.5) * (radius + 0.5)

spawn_locations = []

content = ""
content += "#include \"grid.h\"\n"
content += "\n"
content += "static const int N = NORMAL;\n"
content += "static const int S = SPAWN;\n"
content += "static const int O = OBSTACLE;\n"
content += "\n"
content += "#define ARENA_RADIUS_VALUE " + str(radius) + "\n"
content += "#define GRID_SIZE_VALUE " + str(grid_size) + "\n"

grid_content = "const int GRID[GRID_SIZE_VALUE * GRID_SIZE_VALUE] = {\n"
for i in range(grid_size):
    grid_content += "\t"
    for j in range(grid_size):
        if j != 0:
            grid_content += ", "
        d2 = (i-center)*(i-center) + (j-center)*(j-center)
        if d2 < min2 :
            grid_content += "N"
        elif d2 < max2 :
            spawn_locations.append("(Location){.X = "+str(i)+", .Y = "+str(j)+"}")
            grid_content += "S"
        else :
            grid_content += "O"
    grid_content += ",\n"
grid_content += "};\n"

content += "#define SPAWN_LEN_VALUE " + str(len(spawn_locations)) + "\n"
content += "\n"
content += "const int ARENA_RADIUS = ARENA_RADIUS_VALUE;\n"
content += "const int GRID_SIZE = GRID_SIZE_VALUE;\n"
content += "const int SPAWN_LEN = SPAWN_LEN_VALUE;\n"
content += "\n"

content += grid_content
content += "\n"
content += "const Location SPAWN_LOCATIONS[SPAWN_LEN_VALUE] = {\n"
for i in range(len(spawn_locations)):
    if i != 0:
        content += ",\n"
    content += "\t" + spawn_locations[i]
content += "\n};\n"

output_file.write(content)
output_file.close()