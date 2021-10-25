package setup

import (
	"fmt"

	color "github.com/fatih/color"
)

func PrintSetup(appId string) {
	ccmd := color.New(color.FgGreen).Add(color.Bold)

	fmt.Println("======")
	fmt.Println("🚀 Your new app has been created.")
	fmt.Println("")
	fmt.Println("React/Javascript: To create a brand new project in the current dir:")
	ccmd.Printf("$ npx create-react-app --template @usecloudstate/usecloudstate-js . && sed -i 's/USECLOUDSTATE_APP_ID/%s/' src/utils/cloudstate.js\n", appId)
	fmt.Println("")
	fmt.Println("React/Typescript: To create a brand new project in the current dir:")
	ccmd.Printf("$ npx create-react-app --template @usecloudstate/usecloudstate-ts . && sed -i 's/USECLOUDSTATE_APP_ID/%s/' src/utils/cloudstate.ts\n", appId)
	fmt.Println("")
	fmt.Printf("More docs and admin console is available at: https://usecloudstate.io/admin/apps/%s/settings\n", appId)
	fmt.Println("Happy coding! 👋")
	fmt.Println("======")
}
