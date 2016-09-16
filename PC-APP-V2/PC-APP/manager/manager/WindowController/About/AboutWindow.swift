//
//  AboutWindow.swift
//  SWIFTPROTO
//
//  Created by Almighty Kim on 12/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

import WebKit
import Foundation
import Cocoa

@objc
class AboutWindow: NSWindowController,WebPolicyDelegate {
    
    @IBOutlet weak var webView: WebView!
    var isClosed:Bool = false
    
    override init(window: NSWindow?) {
        super.init(window: window)
    }
    
    required init?(coder: NSCoder) {
        super.init(coder: coder)
    }
 
    override func windowDidLoad() {
        super.windowDidLoad()
        
        var str:String = "<div style=\"text-align:center;font-family:Arial;font-size:13px\">Copyright &copy; {YEAR} Sung-Taek, Kim<br><br>PocketCluster {VERSION}<br><br>For more information visit:<br><a href=\"{URL}\">{URL}</a><br><br>or check us out on GitHub:<br><a href=\"{GITHUB_URL}\">{GITHUB_URL}</a></div>"

        str = str.stringByReplacingOccurrencesOfString("{YEAR}", withString: "2015")
        str = str.stringByReplacingOccurrencesOfString("{VERSION}", withString: (NSBundle.mainBundle().infoDictionary?["CFBundleShortVersionString"] as? String)!)
        str = str.stringByReplacingOccurrencesOfString("{URL}", withString:"https://pocketcluster.wordpress.com")
        str = str.stringByReplacingOccurrencesOfString("{GITHUB_URL}", withString:"https://github.com/stkim1/pocketcluster")
        str = str.stringByReplacingOccurrencesOfString("\n", withString:"<br>")

        self.webView.policyDelegate = self
        self.webView.drawsBackground = false
        self.webView.mainFrame .loadHTMLString(str, baseURL: nil)
        self.isClosed = false
    }
    
    func windowWillClose(notification:NSNotification) {
        NSApplication.sharedApplication().delegate?.performSelector("removeOpenWindow:", withObject: self)
        NSApplication.sharedApplication().endSheet(self.window!, returnCode: 0)
        self.isClosed = true
    }
    
    
    func webView(webView: WebView!, decidePolicyForNavigationAction actionInformation: [NSObject : AnyObject]!, request: NSURLRequest!, frame: WebFrame!, decisionListener listener: WebPolicyDecisionListener!) {
        
        if let _ = request.URL?.host {
            NSWorkspace.sharedWorkspace().openURL(request.URL!)
        } else {
            listener.use()
        }
    }
    
    func use() {}
    func download() {}
    func ignore() {}
}