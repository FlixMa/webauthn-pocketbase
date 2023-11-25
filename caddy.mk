# CADDY
# github.com/caddyserver/caddy

# XCADDY
# https://github.com/caddyserver/xcaddy

# MKCERT
# https://github.com/FiloSottile/mkcert

# HOSTCTL
# https://github.com/guumaster/hostctl
# works on Linux, Mac and Windows :)


# Override variables

# where to look for CaddyFile to use at runtime 
CADDY_SRC_FSPATH=$(PWD)
CADDY_SRC_CONFIG_NAME=Caddyfile
CADDY_SRC_CONFIG_FSPATH=$(PWD)
CADDY_SRC_CONFIG=$(CADDY_SRC_CONFIG_FSPATH)/$(CADDY_SRC_CONFIG_NAME)

CADDY_SRC_GO_BIN_NAME=go
CADDY_SRC_GO_BIN_VERSION=1.21

# defaults for simple setup
CADDY_SRC_DOMAIN=localhost
CADDY_SRC_SUBDOMAINS=

# More advanced setup not finished yet
# CADDY_SRC_DOMAIN=localhost.site
# CADDY_SRC_SUBDOMAINS=localhost.$(CADDY_SRC_DOMAIN) www.$(CADDY_SRC_DOMAIN) sub1.$(CADDY_SRC_DOMAIN) 127.0.0.1


### DEPS

# ! NOTE must use master for a few more days: https://github.com/caddyserver/caddy/issues/5750

CADDY_BIN_NAME=caddy
# https://github.com/caddyserver/caddy/releases/tag/v2.7.5
#CADDY_BIN_VERSION=master
CADDY_BIN_VERSION=v2.7.5
CADDY_BIN_WHICH=$(shell which $(CADDY_BIN_NAME))
CADDY_BIN_WHICH_VERSION=$(shell $(CADDY_BIN_NAME) version)

CADDY_XCADDY_BIN_NAME=xcaddy
# https://github.com/caddyserver/xcaddy/releases/tag/v0.3.5
#CADDY_XCADDY_BIN_VERSION=master
CADDY_XCADDY_BIN_VERSION=v0.3.5
CADDY_XCADDY_BIN_WHICH=$(shell which $(CADDY_XCADDY_BIN_NAME))
CADDY_XCADDY_BIN_WHICH_VERSION=$(shell $(CADDY_XCADDY_BIN_NAME) version)

CADDY_MKCERT_BIN_NAME=mkcert
# https://github.com/FiloSottile/mkcert/releases/tag/v1.4.4
CADDY_MKCERT_BIN_VERSION=v1.4.4
CADDY_MKCERT_BIN_WHICH=$(shell which $(CADDY_MKCERT_BIN_NAME))
CADDY_MKCERT_BIN_WHICH_VERSION=$(shell $(CADDY_MKCERT_BIN_NAME) --version)

# https://github.com/kevinburke/hostsfile
CADDY_HOSTSFILE_BIN_NAME=hostsfile
CADDY_HOSTSFILE_BIN_VERSION=latest
CADDY_HOSTSFILE_BIN_WHICH=$(shell which $(CADDY_HOSTSFILE_BIN_NAME))
CADDY_HOSTSFILE_BIN_WHICH_VERSION=$(shell $(CADDY_HOSTSFILE_BIN_NAME) version)

# https://github.com/guumaster/hostctl
# https://github.com/guumaster/hostctl/releases/tag/v1.1.4
CADDY_HOSTSCTL_BIN_NAME=hostctl
CADDY_HOSTSCTL_BIN_VERSION=v1.1.4
CADDY_HOSTSCTL_BIN_WHICH=$(shell which $(CADDY_HOSTSCTL_BIN_NAME))
CADDY_HOSTSCTL_BIN_WHICH_VERSION=$(shell $(CADDY_HOSTSCTL_BIN_NAME) --version)



# Computed variables
# PERFECT :) https://www.systutorials.com/how-to-get-a-makefiles-directory-for-including-other-makefiles/
_CADDY_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
_CADDY_TEMPLATES_SOURCE=$(_CADDY_SELF_DIR)templates/caddy
_CADDY_TEMPLATES_TARGET=$(PWD)/.templates/caddy



## caddy print, outputs all variables needed to run caddy
caddy-print:
	@echo ""
	@echo "--- CADDY ---"
	@echo ""

caddy-print-src:
	@echo ""
	@echo "Override variables:"
	@echo "CADDY_SRC_GO_BIN_NAME:             $(CADDY_SRC_GO_BIN_NAME)"
	@echo "CADDY_SRC_GO_BIN_VERSION:          $(CADDY_SRC_GO_BIN_VERSION)"
	@echo ""
	@echo "CADDY_SRC_FSPATH:                  $(CADDY_SRC_FSPATH)"
	@echo ""
	@echo "CADDY_SRC_CONFIG_NAME:             $(CADDY_SRC_CONFIG_NAME)"
	@echo "CADDY_SRC_CONFIG_FSPATH:           $(CADDY_SRC_CONFIG_FSPATH)"
	@echo "CADDY_SRC_CONFIG:                  $(CADDY_SRC_CONFIG)"
	@echo ""
	@echo "CADDY_SRC_DOMAIN:                  $(CADDY_SRC_DOMAIN)"
	@echo "CADDY_SRC_SUBDOMAINS:              $(CADDY_SRC_SUBDOMAINS)"
	@echo ""
	

caddy-print-dep:
	@echo ""
	@echo "Deps:"
	@echo "CADDY_BIN_NAME:                    $(CADDY_BIN_NAME)"
	@echo "CADDY_BIN_VERSION:                 $(CADDY_BIN_VERSION)"
	@echo "CADDY_BIN_WHICH:                   $(CADDY_BIN_WHICH)"
	@echo "CADDY_BIN_WHICH_VERSION:           $(CADDY_BIN_WHICH_VERSION)"      


	@echo ""
	@echo "CADDY_XCADDY_BIN_NAME:             $(CADDY_XCADDY_BIN_NAME)"
	@echo "CADDY_XCADDY_BIN_VERSION:          $(CADDY_XCADDY_BIN_VERSION)"
	@echo "CADDY_XCADDY_BIN_WHICH:            $(CADDY_XCADDY_BIN_WHICH)"
	@echo "CADDY_XCADDY_BIN_WHICH_VERSION:    $(CADDY_XCADDY_BIN_WHICH_VERSION)"

	@echo ""
	@echo "CADDY_MKCERT_BIN_NAME:             $(CADDY_MKCERT_BIN_NAME)"
	@echo "CADDY_MKCERT_BIN_VERSION:          $(CADDY_MKCERT_BIN_VERSION)"
	@echo "CADDY_MKCERT_BIN_WHICH:            $(CADDY_MKCERT_BIN_WHICH)"
	@echo "CADDY_MKCERT_BIN_WHICH_VERSION:    $(CADDY_MKCERT_BIN_WHICH_VERSION)"

	@echo ""
	@echo "CADDY_HOSTSFILE_BIN_NAME:          $(CADDY_HOSTSFILE_BIN_NAME)"
	@echo "CADDY_HOSTSFILE_BIN_VERSION:       $(CADDY_HOSTSFILE_BIN_VERSION)"
	@echo "CADDY_HOSTSFILE_BIN_WHICH:         $(CADDY_HOSTSFILE_BIN_WHICH)"
	@echo "CADDY_HOSTSFILE_BIN_WHICH_VERSION: $(CADDY_HOSTSFILE_BIN_WHICH_VERSION)"

	@echo ""
	@echo "CADDY_HOSTSCTL_BIN_NAME:           $(CADDY_HOSTSCTL_BIN_NAME)"
	@echo "CADDY_HOSTSCTL_BIN_VERSION:        $(CADDY_HOSTSCTL_BIN_VERSION)"
	@echo "CADDY_HOSTSCTL_BIN_WHICH:          $(CADDY_HOSTSCTL_BIN_WHICH)"
	@echo "CADDY_HOSTSCTL_BIN_WHICH_VERSION:  $(CADDY_HOSTSCTL_BIN_WHICH_VERSION)"
	@echo ""
	

## caddy dep installs the caddy and mkcert binary to the go bin
## cand copies the templates up into your templates working directory
# Useful where you want to grab them and customise.

## installs go version that caddy may need
caddy-deo-go:
	$(MAKE) GO_SRC_GO_BIN_NAME=go GO_SRC_GO_BIN_VERSION=1.20.1 go-dep-go

## installs go tools i use.
caddy-dep:

	@echo ""
	@echo "installing caddy tool"
	$(CADDY_SRC_GO_BIN_NAME) install -ldflags="-X main.version=$(CADDY_BIN_VERSION)" github.com/caddyserver/caddy/v2/cmd/caddy@$(CADDY_BIN_VERSION)
	@echo ""

	@echo ""
	@echo "installing xcaddy tool"
	# toggled off for now as the Team as a QUIC issue.
	$(CADDY_SRC_GO_BIN_NAME) install -ldflags="-X main.version=$(CADDY_XCADDY_BIN_VERSION)" github.com/caddyserver/xcaddy/cmd/xcaddy@$(CADDY_XCADDY_BIN_VERSION)
	@echo ""

	@echo ""
	@echo "installing mkcert tool"
	$(CADDY_SRC_GO_BIN_NAME) install -ldflags="-X main.version=$(CADDY_MKCERT_BIN_VERSION)" filippo.io/mkcert@$(CADDY_MKCERT_BIN_VERSION)
	@echo ""

	@echo ""
	@echo "installing hostsfile tool"
	# https://github.com/kevinburke/hostsfile
	$(CADDY_SRC_GO_BIN_NAME) install -ldflags="-X main.version=$(CADDY_HOSTSFILE_BIN_VERSION)" github.com/kevinburke/hostsfile@$(CADDY_HOSTSFILE_BIN_VERSION)
	@echo ""

	@echo ""
	#https://github.com/guumaster/hostctl
	@echo "installing hostsctl tool"
	$(CADDY_SRC_GO_BIN_NAME) install -ldflags="-X main.version=$(CADDY_HOSTSCTL_BIN_VERSION)" github.com/guumaster/hostctl/cmd/hostctl@$(CADDY_HOSTSCTL_BIN_VERSION)
	@echo ""

### CADDY APP CONFIG

CADDY_APP_CONFIG_FSPATH="$(HOME)/Library/Application Support/Caddy"
CADDY_APP_CONFIG_AUTOSAVE_FSPATH="$(HOME)/Library/Application Support/Caddy/autosave.json"

caddy-config-ls:
	ls -al $(CADDY_APP_CONFIG_FSPATH)
caddy-config-open:
	open $(CADDY_APP_CONFIG_FSPATH)
caddy-config-autosave-code:
	# opens the resumed config
	# mac only
	code  $(CADDY_APP_CONFIG_AUTOSAVE_FSPATH)
caddy-config-autosave-del:
	rm -f $(CADDY_APP_CONFIG_AUTOSAVE_FSPATH)

### TEMPLATES

## prints the templates 
caddy-templates-print:
	@echo ""
	@echo "- templates:"
	@echo "_CADDY_SELF_DIR:                   $(_CADDY_SELF_DIR)"
	@echo "_CADDY_TEMPLATES_SOURCE:           $(_CADDY_TEMPLATES_SOURCE)"
	@echo "_CADDY_TEMPLATES_TARGET:           $(_CADDY_TEMPLATES_TARGET)"
	@echo ""

caddy-templates-ls:
	@echo ""
	@echo "listing templates ...""
	cd $(_CADDY_TEMPLATES_SOURCE) && ls -al

## installs the caddy templates into your project
caddy-templates-dep: caddy-templates-print
	@echo ""
	@echo "-- installing caddy templates to your project...."
	@echo ""
	mkdir -p $(_CADDY_TEMPLATES_TARGET)
	cp -r $(_CADDY_TEMPLATES_SOURCE)/* $(_CADDY_TEMPLATES_TARGET)
	@echo installed caddy templates  at : $(_CADDY_TEMPLATES_TARGET)
	@echo ""

caddy-mkcert-run-print:
	@echo ""
	@echo "-- mkcert settings"
	@echo "CADDY_SRC_FSPATH:        $(CADDY_SRC_FSPATH)"
	@echo "CADDY_SRC_CONFIG_NAME:   $(CADDY_SRC_CONFIG_NAME)"
	@echo "CADDY_SRC_DOMAIN:        $(CADDY_SRC_DOMAIN)"
	@echo ""

## caddy mkcert installs the certs for browsers to run localhost
caddy-mkcert-run: caddy-mkcert-run-print
	@echo ""
	@echo "-- mkcert is installing certs ..."

	cd $(CADDY_SRC_FSPATH) && $(CADDY_MKCERT_BIN_NAME) -install

	cd $(CADDY_SRC_FSPATH) && $(CADDY_MKCERT_BIN_NAME) $(CADDY_SRC_DOMAIN)

	# TODO finsihed compelx setup..
	# https://www.derpytools.com/how-to-set-up-custom-domain-https-http-3-locally-using-hostsfile-mkcert-caddy/
	# mkcert derpycoder derpycoder.site "*.derpycoder.site" localhost 127.0.0.1 ::1
	#cd $(CADDY_SRC_FSPATH) && $(CADDY_MKCERT_BIN_NAME) $(CADDY_SRC_DOMAIN) "$(CADDY_SRC_DOMAIN)" localhost 127.0.0.1 ::1
	
	@echo ""
	@echo "-- mkcert installed the certs at : $(CADDY_SRC_FSPATH)"
	@echo ""

## mutates hostfile
caddy-hostsfile-run: caddy-mkcert-run-print
	# sudo $(CADDY_HOSTSFILE_BIN_NAME) add derpycoder.site www.derpycoder.site blog.derpycoder.site analytics.derpycoder.site 127.0.0.1
	sudo $(CADDY_HOSTSFILE_BIN_NAME) and $(CADDY_SRC_DOMAIN) $(CADDY_SRC_SUBDOMAINS)

## opens hostfile for checking.
caddy-hostsfile-open:
	code /etc/hosts

## viws hostfile on stdout terminal.
caddy-hostsfile-view:
	cat /etc/hosts



caddy-server-run-print:
	@echo ""
	@echo "-- caddy settings"
	@echo "CADDY_SRC_FSPATH:        $(CADDY_SRC_FSPATH)"
	@echo "CADDY_SRC_CONFIG_NAME:   $(CADDY_SRC_CONFIG_NAME)"
	@echo ""


## caddy fmt runs fixes your caddy file.
caddy-fmt-run:
	cd $(CADDY_SRC_FSPATH) && $(CADDY_BIN_NAME) fmt --overwrite


## caddy without a config for quick stuff
caddy-server-run-browse: caddy-server-run-print
	@echo ""
	@echo "-- caddy browser starting ..."
	@echo ""
	cd $(CADDY_SRC_FSPATH) && $(CADDY_BIN_NAME) file-server --browse

## caddy run runs caddy using your Caddyfile. 
caddy-server-run: caddy-server-run-print 
	@echo ""
	@echo "-- caddy server starting ..."
	@echo ""
	cd $(CADDY_SRC_FSPATH) && $(CADDY_BIN_NAME) run --config $(CADDY_SRC_CONFIG)

	# --environ --resume. DONT use this
	#  --config and --resume flags were used together; ignoring --config and resuming from last configuration  {"autosave_file": "/Users/apple/Library/Application Support/Caddy/autosave.json"}
	#cd $(CADDY_SRC_FSPATH) && $(CADDY_BIN_NAME) run --environ --resume --config $(CADDY_SRC_CONFIG_NAME)
	
	# open https://localhost:8443

## caddy run and watch for config changes	
caddy-server-run-watch: caddy-server-run-print
	@echo ""
	@echo "-- caddy server starting in watch mode ..."
	@echo ""
	cd $(CADDY_SRC_FSPATH) && $(CADDY_BIN_NAME) run --config $(CADDY_SRC_CONFIG_NAME) --watch


