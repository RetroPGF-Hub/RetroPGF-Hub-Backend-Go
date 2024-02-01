#!/bin/bash
# chmod +x generate-module.sh
# ./generate-module.sh product

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <filename>"
  exit 1
fi

FILENAME=$1
FIRST_CAPITALIZED_F=$(echo "${FILENAME}" | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2))}')

# Create directory structure
mkdir -p "modules/${FILENAME}"
cd "modules/${FILENAME}"

# Create files
echo "package ${FILENAME}" > ${FILENAME}Entity.go
echo "package ${FILENAME}" > ${FILENAME}Model.go

# ---------------- HTTP SECTION ---------------- 
mkdir "${FILENAME}HttpHandler"
# echo "package ${FILENAME}handler" > ${FILENAME}Handler/${FILENAME}HttpHandler.go
cat <<EOL > "${FILENAME}HttpHandler/${FILENAME}HttpHandler.go"
package ${FILENAME}httphandler

type (
	${FIRST_CAPITALIZED_F}HttpHandlerService interface {
	}

	${FILENAME}HttpHandler struct {
		${FILENAME}Usecase ${FILENAME}usecase.${FIRST_CAPITALIZED_F}UsecaseService
	}
)

func New${FIRST_CAPITALIZED_F}HttpHandler(${FILENAME}Usecase ${FILENAME}usecase.${FIRST_CAPITALIZED_F}UsecaseService) ${FIRST_CAPITALIZED_F}HttpHandlerService {
	return &${FILENAME}HttpHandler{
		${FILENAME}Usecase: ${FILENAME}Usecase,
	}
}
EOL
# ---------------- HTTP SECTION ---------------- 


# ---------------- USECASE ---------------- 
mkdir "${FILENAME}Usecase"
cat <<EOL > "${FILENAME}Usecase/${FILENAME}Usecase.go"
package ${FILENAME}usecase

type (
	${FIRST_CAPITALIZED_F}UsecaseService interface {

	}

	${FILENAME}Usecase struct {
		${FILENAME}Repo ${FILENAME}repository.${FIRST_CAPITALIZED_F}RepositoryService
	}
)

func New${FIRST_CAPITALIZED_F}Usecase(${FILENAME}Repo ${FILENAME}repository.${FIRST_CAPITALIZED_F}RepositoryService) ${FIRST_CAPITALIZED_F}UsecaseService {
	return &${FILENAME}Usecase{
		${FILENAME}Repo: ${FILENAME}Repo,
	}
}
EOL
# ---------------- USECASE ---------------- 


# ---------------- REPOSITORY ---------------- 
mkdir "${FILENAME}Repository"
cat <<EOL > "${FILENAME}Repository/${FILENAME}Repository.go"
package ${FILENAME}repository

type (
	${FIRST_CAPITALIZED_F}RepositoryService interface {

	}

	${FILENAME}Repository struct {
		db     any
	}
)

func New${FIRST_CAPITALIZED_F}Repository(db any) ${FIRST_CAPITALIZED_F}RepositoryService {
	return &${FILENAME}Repository{
		db:     db,
	}
}
EOL
# ---------------- REPOSITORY ---------------- 


# ---------------- SERVER Module Name ----------------
cd ../..
mkdir -p "server"
touch "server/${FILENAME}.go"

cat <<EOL > "server/${FILENAME}.go"
package server

func (s *server) ${FILENAME}Service() {
	${FILENAME}Repo := ${FILENAME}repository.New${FIRST_CAPITALIZED_F}Repository(s.db)
	${FILENAME}Usecase := ${FILENAME}usecase.New${FIRST_CAPITALIZED_F}Usecase(${FILENAME}Repo)
	${FILENAME}HttpHandler := ${FILENAME}httphandler.New${FIRST_CAPITALIZED_F}HttpHandler(${FILENAME}Usecase)
	
	_ = ${FILENAME}HttpHandler
}

EOL
# ---------------- SERVER Module Name ----------------

echo "Folder and file structure for ${FILENAME} created successfully."