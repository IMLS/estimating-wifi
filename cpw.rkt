#lang racket

(require racket/cmdline)

(define LOGICAL-MATCH-PATTERN #px"[[:space:]+]logical name:[[:space:]+](.*)")

(define (find-ralink los #:usb? [usb? false] #:hash [h (make-immutable-hash)])
  (cond
    [(empty? los) h]
    [(regexp-match #px"[[:space:]*]\\*-usb" (first los))
     (find-ralink (rest los) #:usb? true)]
    [(and usb? (regexp-match #px"[[:space:]*]\\*-.*" (first los)))
     h]
    [usb?
     (define line (regexp-match #px"^\\s+([[:alnum:][:space:]]+):\\s*(.*?)$" (first los)))
     ;;(printf "~s~n" line)
     (find-ralink (rest los) #:usb? true #:hash (hash-set h (second line) (third line)))] 
    [else
     (find-ralink (rest los) #:usb? usb? #:hash h)]))

(define ()
   (define result
     (with-output-to-string
       (thunk
        (system* "/usr/bin/lshw" "-class" "network"))))
   (define lines (regexp-split #px"\n" result))
   (define device-hash (find-ralink lines))
   ;; Now, if it is RAlink, return the interface.
   (printf "~a~n" (hash-ref device-hash "logical name" "NOTFOUND"))
   (if (hash-has-key? device-hash "logical name")
       (exit 0)
       (exit -1))
       )
(define (cpw)
  (command-line
   #:program "configure-pi-wifi"
   #:args ()

   ))
